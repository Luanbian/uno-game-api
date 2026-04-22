// Package action provides flux of possible actions in gameplay
package action

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Luanbian/uno-game-api/internal/hub"
	"github.com/gorilla/websocket"
)

type Action string

const (
	ActionJoin           Action = "join"
	ActionStartGame      Action = "start_game"
	ActionPlayCard       Action = "play_card"
	ActionBuyCard        Action = "buy_card"
	ActionSayUno         Action = "say_uno"
	ActionPunishNoSayUno Action = "punish_no_say_uno"
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

	switch Action(payload.Action) {
	case ActionJoin:
		hub.AddNewPlayer(payload.Nickname, connection)
	default:
		return nil, fmt.Errorf("unknown action: %s", payload.Action)
	}

	return []byte(""), nil
}
