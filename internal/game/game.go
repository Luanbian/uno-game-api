// Package game provides a game state
package game

import (
	"encoding/json"
	"fmt"
	"sync"
)

type GameState struct {
	mutex         sync.Mutex
	Deck          *Deck             `json:"-"`
	Hands         map[string][]Card `json:"hands"`
	DiscardPile   []Card            `json:"discard_pile"`
	Players       []string          `json:"players"`
	CurrentPlayer int               `json:"current_player"`
	LastPlayer    int               `json:"last_player"`
	SaidUno       map[string]bool   `json:"said_uno"`
	Winner        string            `json:"winner,omitempty"`
}

var (
	gameState *GameState
	mutex     sync.RWMutex
)

func GameStateToJSON(gs *GameState) ([]byte, error) {
	result, err := json.Marshal(gs)
	if err != nil {
		return nil, fmt.Errorf("converting game state to json: %w", err)
	}

	return result, nil
}

func SetCurrentGameState(gs *GameState) {
	mutex.Lock()
	defer mutex.Unlock()

	gameState = gs
}

func GetCurrentGameState() (*GameState, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	if gameState == nil {
		return nil, fmt.Errorf("game not started")
	}

	return gameState, nil
}
