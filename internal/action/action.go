// Package action provides flux of possible actions in gameplay
package action

import (
	"encoding/json"
	"fmt"

	"github.com/Luanbian/uno-game-api/internal/game"
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
		return nil, fmt.Errorf("converting message to Payload: %w", err)
	}

	switch Action(payload.Action) {
	case ActionJoin:
		hub.AddNewPlayer(payload.Nickname, connection)
		return []byte("Player adicionado"), nil
	case ActionStartGame:
		return startGame()
	default:
		return nil, fmt.Errorf("unknown action: %s", payload.Action)
	}
}

func startGame() ([]byte, error) {
	players := hub.GetPlayers()
	gameState, err := game.StartGame(players)
	if err != nil {
		return nil, err
	}

	result, err := game.GameStateToJSON(gameState)
	if err != nil {
		return nil, err
	}

	return result, nil
}
