package utils

import (
	db "chess_server/database"
	"chess_server/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var Players = make(map[uint]*websocket.Conn)
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

type MoveMessage struct {
	GameID uint   `json:"game_id"`
	Move   string `json:"move"`
}

func HandleConnection(playerId uint, w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrader.Upgrade(w, r, nil)
	defer func() {
		delete(Players, playerId)
		ws.Close()
	}()
	if err != nil {
		return
	}
	defer ws.Close()

	Players[playerId] = ws

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			delete(Players, playerId)
			break
		}

		var moveMsg MoveMessage
		if err := json.Unmarshal(msg, &moveMsg); err != nil {
			fmt.Println("Invalid message:", err)
			continue
		}

		var game models.Game
		if err := db.DB.First(&game, moveMsg.GameID).Error; err != nil {
			fmt.Println("Game not found:", err)
			continue
		}

		var moves []string
		if len(game.Moves) > 0 {
			if err := json.Unmarshal(game.Moves, &moves); err != nil {
				fmt.Println("Failed to unmarshal moves:", err)
				moves = []string{}
			}
		}
		moves = append(moves, moveMsg.Move)

		newMoves, _ := json.Marshal(moves)
		game.Moves = newMoves

		if err := db.DB.Save(&game).Error; err != nil {
			fmt.Println("Failed to save move:", err)
		}

		var opponentID uint
		if game.Player1ID == playerId {
			opponentID = game.Player2ID
		} else if game.Player2ID == playerId {
			opponentID = game.Player1ID
		} else {
			fmt.Println("Player not part of this game")
			continue
		}

		if opponentWS, ok := Players[opponentID]; ok {
			if err := opponentWS.WriteMessage(websocket.TextMessage, msg); err != nil {
				opponentWS.Close()
				delete(Players, opponentID)
			}
		}
	}
}
