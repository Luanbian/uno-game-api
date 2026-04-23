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
	Action   string    `json:"action"`
	Nickname string    `json:"nickname"`
	Card     game.Card `json:"card"`
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
	case ActionPlayCard:
		return playCard(payload.Nickname, payload.Card)
	case ActionBuyCard:
		return buyCard(payload.Nickname)
	case ActionSayUno:
		return sayUno(payload.Nickname)
	case ActionPunishNoSayUno:
		return punishNoSayUno()
	default:
		return nil, fmt.Errorf("unknown action: %s", payload.Action)
	}
}

func hasWinner() (bool, error) {
	gameState, err := game.GetCurrentGameState()
	if err != nil {
		return false, fmt.Errorf("getting current game state: %w", err)
	}

	if gameState.Winner != "" {
		return true, nil
	}

	return false, nil
}

func startGame() ([]byte, error) {
	players := hub.GetPlayers()

	gameState, err := game.GetCurrentGameState()
	if gameState != nil && err == nil {
		err = game.Rematch(players)
		if err != nil {
			return nil, err
		}
	}

	gameState, err = game.StartGame(players)
	if err != nil {
		return nil, err
	}

	game.SetCurrentGameState(gameState)

	result, err := game.GameStateToJSON(gameState)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func playCard(nickname string, card game.Card) ([]byte, error) {
	playerWin, err := hasWinner()
	if err != nil {
		return nil, err
	}
	if playerWin {
		return nil, fmt.Errorf("game already has a winner")
	}

	gameState, err := game.PlayCard(nickname, card)
	if err != nil {
		return nil, err
	}

	result, err := game.GameStateToJSON(gameState)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func buyCard(nickname string) ([]byte, error) {
	playerWin, err := hasWinner()
	if err != nil {
		return nil, err
	}
	if playerWin {
		return nil, fmt.Errorf("game already has a winner")
	}

	gameState, err := game.BuyCard(nickname)
	if err != nil {
		return nil, err
	}

	result, err := game.GameStateToJSON(gameState)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func sayUno(nickname string) ([]byte, error) {
	playerWin, err := hasWinner()
	if err != nil {
		return nil, err
	}
	if playerWin {
		return nil, fmt.Errorf("game already has a winner")
	}

	gameState, err := game.SayUno(nickname)
	if err != nil {
		return nil, err
	}

	result, err := game.GameStateToJSON(gameState)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func punishNoSayUno() ([]byte, error) {
	playerWin, err := hasWinner()
	if err != nil {
		return nil, err
	}
	if playerWin {
		return nil, fmt.Errorf("game already has a winner")
	}

	gameState, err := game.PunishNoSayUno()
	if err != nil {
		return nil, err
	}

	result, err := game.GameStateToJSON(gameState)
	if err != nil {
		return nil, err
	}

	return result, nil
}
