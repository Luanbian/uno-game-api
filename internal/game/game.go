// Package game provides a game state
package game

import (
	"encoding/json"
	"fmt"
)

type GameState struct {
	Deck          *Deck             `json:"-"`
	Hands         map[string][]Card `json:"hands"`
	DiscardPile   []Card            `json:"discard_pile"`
	Players       []string          `json:"players"`
	CurrentPlayer int               `json:"current_player"`
}

func GameStateToJSON(gameState *GameState) ([]byte, error) {
	result, err := json.Marshal(gameState)
	if err != nil {
		return nil, fmt.Errorf("converting game state to json: %w", err)
	}

	return result, nil
}
