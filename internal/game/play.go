// Package game play controls each play of each player
package game

import "fmt"

func PlayCard(nickname string, card Card) (*GameState, error) {
	currentGame, err := GetCurrentGameState()
	if err != nil {
		return nil, fmt.Errorf("playing card: %w: ", err)
	}
	currentGame.mutex.Lock()
	defer currentGame.mutex.Unlock()

	if isColorSelectPending(currentGame) {
		return nil, fmt.Errorf("waiting for color selection")
	}

	err = isYourTurn(nickname, currentGame)
	if err != nil {
		return nil, err
	}

	err = isCardValid(card, currentGame)
	if err != nil {
		return nil, err
	}

	err = removeCardFromHand(nickname, card, currentGame)
	if err != nil {
		return nil, err
	}

	addCardInTopOfPile(card, currentGame)

	hasWinner, err := playerWin(nickname, currentGame)
	if err != nil {
		return nil, err
	}
	if hasWinner {
		return currentGame, nil
	}

	playerIndex := currentGame.CurrentPlayer

	err = ApplyCardEffect(card, currentGame)
	if err != nil {
		return nil, err
	}

	if card.Color == None {
		currentGame.LastPlayer = playerIndex
		return currentGame, nil
	}

	nextTurn(currentGame)

	return currentGame, nil
}

func isYourTurn(nickname string, gs *GameState) error {
	currentPlayer := gs.Players[gs.CurrentPlayer]
	if currentPlayer != nickname {
		return fmt.Errorf("not your turn")
	}

	return nil
}

func isCardValid(playedCard Card, gs *GameState) error {
	if !canPlayCard(playedCard, gs) {
		return fmt.Errorf("invalid card")
	}
	return nil
}

func removeCardFromHand(nickname string, card Card, gs *GameState) error {
	hand, ok := gs.Hands[nickname]
	if !ok {
		return fmt.Errorf("removing card from non existing hand: %s", nickname)
	}

	foundIndex := -1
	for i, c := range hand {
		if c == card {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return fmt.Errorf("not found card to remove")
	}

	hand[foundIndex] = hand[len(hand)-1]
	hand = hand[:len(hand)-1]

	gs.Hands[nickname] = hand

	return nil
}

func addCardInTopOfPile(card Card, gs *GameState) {
	gs.DiscardPile = append(gs.DiscardPile, card)
}

func nextTurn(gs *GameState) {
	currentPlayer := gs.CurrentPlayer
	nextPlayer := (currentPlayer + gs.Direction + len(gs.Players)) % len(gs.Players)
	gs.LastPlayer = currentPlayer
	gs.CurrentPlayer = nextPlayer
}

func playerWin(nickname string, gs *GameState) (bool, error) {
	hand, ok := gs.Hands[nickname]
	if !ok {
		return false, fmt.Errorf("checking win for non existing hand: %s", nickname)
	}

	if len(hand) == 0 {
		gs.Winner = nickname
		return true, nil
	}

	return false, nil
}

func isColorSelectPending(gs *GameState) bool {
	top := gs.DiscardPile[len(gs.DiscardPile)-1]
	return top.Color == None
}
