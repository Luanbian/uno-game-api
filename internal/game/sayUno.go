// Package game say uno provide a feature to say uno
package game

import "fmt"

func SayUno(nickname string) (*GameState, error) {
	currentGame, err := GetCurrentGameState()
	if err != nil {
		return nil, fmt.Errorf("saying uno: %w: ", err)
	}
	currentGame.mutex.Lock()
	defer currentGame.mutex.Unlock()

	err = areYouLastPlayer(nickname, currentGame)
	if err != nil {
		return nil, err
	}

	err = hasOneCard(nickname, currentGame)
	if err != nil {
		return nil, err
	}

	err = alreadySaidUno(nickname, currentGame)
	if err != nil {
		return nil, err
	}

	currentGame.SaidUno[nickname] = true

	return currentGame, nil
}

func areYouLastPlayer(nickname string, gs *GameState) error {
	if gs.LastPlayer == -1 {
		return fmt.Errorf("no last player")
	}
	if nickname != gs.Players[gs.LastPlayer] {
		return fmt.Errorf("you are not the last player")
	}

	return nil
}

func hasOneCard(nickname string, gs *GameState) error {
	hand, ok := gs.Hands[nickname]
	if !ok {
		return fmt.Errorf("checking hand of non existing player: %s", nickname)
	}

	if len(hand) != 1 {
		return fmt.Errorf("you don't have one card left")
	}

	return nil
}

func alreadySaidUno(nickname string, gs *GameState) error {
	if gs.SaidUno[nickname] {
		return fmt.Errorf("you already said uno")
	}

	return nil
}

func resetSaidUno(nickname string, gs *GameState) error {
	if _, ok := gs.SaidUno[nickname]; !ok {
		return fmt.Errorf("resetting said uno for non existing player: %s", nickname)
	}
	gs.SaidUno[nickname] = false

	return nil
}
