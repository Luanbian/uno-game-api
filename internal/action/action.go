// Package action provides flux of possible actions in gameplay
package action

import (
	"encoding/json"
	"log"

	"github.com/Luanbian/uno-game-api/internal/hub"
	"github.com/gorilla/websocket"
)

type Payload struct {
	Action   string `json:"action"`
	Nickname string `json:"nickname"`
}

func Handler(message []byte, connection *websocket.Conn) ([]byte, error) {
	var payload Payload

	err := json.Unmarshal(message, &payload)
	if err != nil {
		log.Println("Error to convert message in Payload: ", err)
		return nil, err
	}

	if payload.Action == "join" {
		hub.AddNewPlayer(payload.Nickname, connection)
	}

	return []byte("tudo ok agora"), nil
}
