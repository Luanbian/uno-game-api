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
	Direction     int               `json:"direction"`
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

func (gs *GameState) FilterForPlayer(nickname string) *GameState {
	filteredHands := make(map[string][]Card, len(gs.Hands))

	for player, hand := range gs.Hands {
		if player == nickname {
			filteredHands[player] = hand
		} else {
			filteredHands[player] = make([]Card, len(hand))
		}
	}

	gameState := &GameState{
		Deck:          gs.Deck,
		Hands:         filteredHands,
		DiscardPile:   gs.DiscardPile,
		Players:       gs.Players,
		CurrentPlayer: gs.CurrentPlayer,
		LastPlayer:    gs.LastPlayer,
		SaidUno:       gs.SaidUno,
		Winner:        gs.Winner,
		Direction:     gs.Direction,
	}

	return gameState
}

func RemovePlayerFromGameState(nickname string) error {
	gs, err := GetCurrentGameState()
	if err != nil {
		return fmt.Errorf("removing player from game: %w", err)
	}

	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	hand, ok := gs.Hands[nickname]
	if !ok {
		return fmt.Errorf("player %s not found in game state", nickname)
	}

	gs.Deck.cards = append(gs.Deck.cards, hand...)
	gs.Deck.shuffle()

	delete(gs.Hands, nickname)
	delete(gs.SaidUno, nickname)

	for i, player := range gs.Players {
		if player == nickname {
			gs.Players = append(gs.Players[:i], gs.Players[i+1:]...)
			if gs.CurrentPlayer == i {
				gs.CurrentPlayer = i % len(gs.Players)
			} else if gs.CurrentPlayer > i {
				gs.CurrentPlayer--
			}
			break
		}
	}

	return nil
}
