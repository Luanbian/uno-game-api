// Package game punish player who didn't say uno when they had one card left
package game

import "fmt"

func PunishNoSayUno() (*GameState, error) {
	currentGame, err := GetCurrentGameState()
	if err != nil {
		return nil, fmt.Errorf("punishing no say uno: %w: ", err)
	}
	currentGame.mutex.Lock()
	defer currentGame.mutex.Unlock()

	err = isLastPlayerHasOneCard(currentGame)
	if err != nil {
		return nil, err
	}

	err = isLastPlayerAlreadySaidUno(currentGame)
	if err != nil {
		return nil, err
	}

	lastPlayer := currentGame.Players[currentGame.LastPlayer]

	var topCard Card
	for range 2 {
		topCard, err = currentGame.Deck.pickUpCard()
		if err != nil {
			return nil, err
		}

		err = addCardInHand(lastPlayer, topCard, currentGame)
		if err != nil {
			return nil, err
		}
	}

	err = resetSaidUno(lastPlayer, currentGame)
	if err != nil {
		return nil, err
	}

	return currentGame, nil
}

func isLastPlayerHasOneCard(gs *GameState) error {
	if gs.LastPlayer == -1 {
		return fmt.Errorf("no last player")
	}
	lastPlayer := gs.Players[gs.LastPlayer]

	hand, ok := gs.Hands[lastPlayer]
	if !ok {
		return fmt.Errorf("checking hand of non existing player: %s", lastPlayer)
	}

	if len(hand) != 1 {
		return fmt.Errorf("last player doesn't have one card left")
	}

	return nil
}

func isLastPlayerAlreadySaidUno(gs *GameState) error {
	if gs.LastPlayer == -1 {
		return fmt.Errorf("no last player")
	}

	lastPlayer := gs.Players[gs.LastPlayer]
	if gs.SaidUno[lastPlayer] {
		return fmt.Errorf("last player said uno")
	}

	return nil
}
