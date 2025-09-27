package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var Players = make(map[uint]*websocket.Conn)
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
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

		//TODO: instead of echo back the message add the game logic.
		if err := Players[playerId].WriteMessage(websocket.TextMessage, msg); err != nil {
			Players[playerId].Close()
			delete(Players, playerId)
		}
	}
}
