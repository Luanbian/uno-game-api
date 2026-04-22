// Package ws provider websocket configuration
package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/Luanbian/uno-game-api/internal/action"
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
		msgtype, message, err := connection.ReadMessage()
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
	delete(wsServer.clients, connection)
	wsServer.mutex.Unlock()

	_ = connection.Close()
}

func (wsServer *WsServer) handleMessage(message []byte, connection *websocket.Conn) {
	response, err := action.Handler(message, connection)
	if err != nil {
		log.Println("Action error: ", err)
		return
	}
	if len(response) == 0 {
		return
	}

	wsServer.WriteMessage(response)
}

func (wsServer *WsServer) WriteMessage(message []byte) {
	wsServer.mutex.RLock()
	defer wsServer.mutex.RUnlock()

	for conn := range wsServer.clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Write error: ", err)
			continue
		}
	}
}
