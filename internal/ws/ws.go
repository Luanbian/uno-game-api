// Package ws provider websocket configuration
package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const Route = "GET /ws"

type WsServer struct {
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

	wsServer.clients[connection] = true

	for {
		msgtype, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("Message error: ", err)
			break
		}
		if msgtype == websocket.CloseMessage {
			break
		}

		go wsServer.handleMessage(message)
	}

	delete(wsServer.clients, connection)

	_ = connection.Close()
}

func (wsServer *WsServer) handleMessage(message []byte) {
	log.Println("Received: ", string(message))
	wsServer.WriteMessage(message)
}

func (wsServer *WsServer) WriteMessage(message []byte) {
	for conn := range wsServer.clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Write error: ", err)
			break
		}
	}
}
