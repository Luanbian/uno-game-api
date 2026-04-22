// Package ws provider websocket configuration
package ws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const Route = "GET /ws"

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

	for {
		msgtype, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("Message error: ", err)
			break
		}
		if msgtype == websocket.CloseMessage {
			break
		}

		if err := connection.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Write message error: ", err)
			break
		}

		go messageHandler(message)
	}

	_ = connection.Close()
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}
