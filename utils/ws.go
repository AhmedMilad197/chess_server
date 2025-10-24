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

type Message struct {
	GameID int             `json:"game_id"`
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
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

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			fmt.Println("Invalid message:", err)
			continue
		}
		var game models.Game
		if err := db.DB.First(&game, message.GameID).Error; err != nil {
			fmt.Println("Game not found:", err)
			continue
		}
		var moveData struct {
			From string `json:"from"`
			To   string `json:"to"`
		}
		if message.Type == "move" {
			if err := json.Unmarshal(message.Data, &moveData); err != nil {
				fmt.Println("Error unmarshaling move data:", err)
				return
			}

			var moves []string
			if len(game.Moves) > 0 {
				if err := json.Unmarshal(game.Moves, &moves); err != nil {
					fmt.Println("Failed to unmarshal moves:", err)
					moves = []string{}
				}
			}
			moves = append(moves, moveData.To)
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
}
