package utils

import (
	"chess_server/database"
	"chess_server/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

var Ctx context.Context

type NotificationMessage struct {
	Type     string `json:"type"`
	GameId   int    `json:"game_id"`
	Opponent Player `json:"opponent"`
	IsBlack  bool   `json:"is_black"`
}

func InitGame() {
	Ctx = context.Background()
}

func createGame(p1, p2 Player) {
	game := models.Game{
		Player1ID:  p1.UserID,
		Player2ID:  p2.UserID,
		GameTypeID: p1.GameTypeID,
		Status:     "ongoing",
	}
	db.DB.Create(&game)
	turn := rand.Intn(2)
	isBlack := true
	if turn == 1 {
		isBlack = false
	}
	startGame := "start_game"
	player1Notification := NotificationMessage{
		Type:     startGame,
		GameId:   int(game.ID),
		Opponent: p2,
		IsBlack:  isBlack,
	}
	player2Notification := NotificationMessage{
		Type:     startGame,
		GameId:   int(game.ID),
		Opponent: p1,
		IsBlack:  !isBlack,
	}
	player1Data, _ := json.Marshal(&player1Notification)
	player2Data, _ := json.Marshal(&player2Notification)

	if err := Players[p1.UserID].WriteMessage(websocket.TextMessage, []byte(player1Data)); err != nil {
		Players[p1.UserID].Close()
		delete(Players, p1.UserID)
	}
	if err := Players[p2.UserID].WriteMessage(websocket.TextMessage, []byte(player2Data)); err != nil {
		Players[p1.UserID].Close()
		delete(Players, p1.UserID)
	}

	fmt.Printf("Game created: %d vs %d\n", p1.UserID, p2.UserID)
}

type Player struct {
	UserID               uint
	GameTypeID           uint
	Rating               int
	LowerBoundRatingDiff int
	UpperBoundRatingDiff int
}

func mutualFit(p1, p2 Player) bool {
	return p2.Rating >= p1.Rating-p1.LowerBoundRatingDiff &&
		p2.Rating <= p1.Rating+p1.UpperBoundRatingDiff &&
		p1.Rating >= p2.Rating-p2.LowerBoundRatingDiff &&
		p1.Rating <= p2.Rating+p2.UpperBoundRatingDiff
}

func EnqueuePlayer(userId uint, gameTypeId int) {
	var user models.User
	db.DB.Preload("Ratings.GameType").Preload("Setting").First(&user, userId)

	var playerRating int
	for _, rating := range user.Ratings {
		if rating.GameTypeID == uint(gameTypeId) {
			playerRating = rating.Rating
		}
	}

	player := Player{
		UserID:               user.ID,
		GameTypeID:           uint(gameTypeId),
		Rating:               playerRating,
		LowerBoundRatingDiff: int(user.Setting.LowerBoundPlayerRatingDiff),
		UpperBoundRatingDiff: int(user.Setting.UpperBoundPlayerRatingDiff),
	}

	serialized, err := json.Marshal(player)
	if err != nil {
		fmt.Println("Error marshaling player:", err)
		return
	}
	serializedStr := string(serialized)

	exists, err := RDB.SIsMember(Ctx, "players_q_set", serializedStr).Result()
	if err != nil {
		fmt.Println("Error checking set:", err)
		return
	}

	if exists {
		fmt.Println("Player already in queue, skipping")
		return
	}

	pipe := RDB.TxPipeline()
	pipe.SAdd(Ctx, "players_q_set", serializedStr)
	pipe.RPush(Ctx, "players_q", serializedStr)
	_, err = pipe.Exec(Ctx)
	if err != nil {
		fmt.Println("Error enqueuing player:", err)
		return
	}

	fmt.Println("Player enqueued")
}

func MatchmakingWorker() {
	for {
		players, err := RDB.LRange(Ctx, "players_q", 0, -1).Result()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if len(players) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		for _, playerRaw := range players {
			matched := MatchPlayer(playerRaw)

			if matched {
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func MatchPlayer(playerRaw string) bool {
	err := RDB.Watch(Ctx, func(tx *redis.Tx) error {
		players, err := tx.LRange(Ctx, "players_q", 0, -1).Result()
		if err != nil {
			return err
		}

		var p Player
		if err := json.Unmarshal([]byte(playerRaw), &p); err != nil {
			return err
		}

		for _, raw := range players {
			var candidate Player
			if err := json.Unmarshal([]byte(raw), &candidate); err != nil {
				continue
			}

			if candidate.UserID == p.UserID {
				continue
			}
			if candidate.GameTypeID == p.GameTypeID && mutualFit(p, candidate) {
				pipe := tx.TxPipeline()
				pipe.LRem(Ctx, "players_q", 1, raw)
				pipe.LRem(Ctx, "players_q", 1, playerRaw)
				pipe.SRem(Ctx, "players_q_set", raw)
				pipe.SRem(Ctx, "players_q_set", playerRaw)
				_, err := pipe.Exec(Ctx)
				if err != nil {
					return err
				}

				fmt.Printf("Matched Player %d with Player %d in gameType %d\n",
					p.UserID, candidate.UserID, p.GameTypeID)

				/*
				* create game in DB
				* TODO: notify players
				 */
				createGame(p, candidate)

				return nil
			}
		}

		return nil
	}, "players_q")

	if err != nil {
		fmt.Println("Matchmaking transaction failed:", err)
		return false
	}

	return true
}
