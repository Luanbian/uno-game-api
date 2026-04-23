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
	topOfPile := gs.DiscardPile[len(gs.DiscardPile)-1]
	if playedCard.Color == None {
		return nil
	}

	sameColor := playedCard.Color == topOfPile.Color
	sameType := playedCard.Type == topOfPile.Type && playedCard.Type != Number
	sameValue := playedCard.Type == Number && topOfPile.Type == Number && playedCard.Value ==
		topOfPile.Value

	if !sameColor && !sameType && !sameValue {
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
	nextPlayer := (currentPlayer + 1) % len(gs.Players)

	gs.CurrentPlayer = nextPlayer
}
