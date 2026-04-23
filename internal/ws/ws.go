// Package ws provider websocket configuration
package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Luanbian/uno-game-api/internal/action"
	"github.com/Luanbian/uno-game-api/internal/game"
	"github.com/Luanbian/uno-game-api/internal/hub"
	"github.com/gorilla/websocket"
)

const Route = "GET /ws"

type WsServer struct {
	mutex   sync.RWMutex
	clients map[*websocket.Conn]bool
}

var wsServer = &WsServer{
	clients: make(map[*websocket.Conn]bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(res http.ResponseWriter, req *http.Request) {
	connection, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println("Connection error: ", err)
		return
	}

	wsServer.mutex.Lock()
	wsServer.clients[connection] = true
	wsServer.mutex.Unlock()

	for {
		var msgtype int
		var message []byte
		msgtype, message, err = connection.ReadMessage()
		if err != nil {
			log.Println("Message error: ", err)
			break
		}
		if msgtype == websocket.CloseMessage {
			break
		}

		go wsServer.handleMessage(message, connection)
	}

	wsServer.mutex.Lock()
	defer wsServer.mutex.Unlock()
	defer func() {
		if err = connection.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}()

	nickname, err := hub.RemovePlayerByConn(connection)
	if err != nil {
		log.Println("Error removing player: ", err)
		delete(wsServer.clients, connection)
		return
	}
	err = game.RemovePlayerFromGameState(nickname)
	if err != nil {
		log.Println("Error removing player from game state: ", err)
	}
	delete(wsServer.clients, connection)
}

func (wsServer *WsServer) handleMessage(message []byte, connection *websocket.Conn) {
	response, err := action.Handler(message, connection)
	if err != nil {
		errMsg, _ := json.Marshal(err.Error())
		buffer := []byte{}
		buffer = fmt.Appendf(buffer, `{"error":%s}`, errMsg)

		if err = connection.WriteMessage(
			websocket.TextMessage,
			buffer,
		); err != nil {
			log.Println("Write error: ", err)
		}
		return
	}

	if response != nil {
		wsServer.BroadcastGameState(response)
	}
}

func (wsServer *WsServer) BroadcastGameState(gs *game.GameState) {
	connections := hub.GetPlayerConnections()
	for nickname, conn := range connections {
		filtered := gs.FilterForPlayer(nickname)
		result, err := game.GameStateToJSON(filtered)
		if err != nil {
			continue
		}
		err = conn.WriteMessage(websocket.TextMessage, result)
		if err != nil {
			log.Printf("Error broadcasting to %s: %v", nickname, err)
		}
	}
}
