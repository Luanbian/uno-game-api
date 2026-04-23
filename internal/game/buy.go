// Package game buy provides a logic to buy a card from deck
package game

import "fmt"

func BuyCard(nickname string) (*GameState, error) {
	currentGame, err := GetCurrentGameState()
	if err != nil {
		return nil, fmt.Errorf("buying card: %w: ", err)
	}
	currentGame.mutex.Lock()
	defer currentGame.mutex.Unlock()

	err = isYourTurn(nickname, currentGame)
	if err != nil {
		return nil, err
	}

	topCard, err := currentGame.Deck.pickUpCard()
	if err != nil {
		return nil, err
	}

	err = addCardInHand(nickname, topCard, currentGame)
	if err != nil {
		return nil, err
	}

	if !canPlayCard(topCard, currentGame) {
		nextTurn(currentGame)
	}

	return currentGame, nil
}

func addCardInHand(nickname string, card Card, gs *GameState) error {
	hand, ok := gs.Hands[nickname]
	if !ok {
		return fmt.Errorf("adding card from non existing hand: %s", nickname)
	}

	gs.Hands[nickname] = append(hand, card)
	return nil
}

func canPlayCard(card Card, gs *GameState) bool {
	if card.Color == None {
		return true
	}
	top := gs.DiscardPile[len(gs.DiscardPile)-1]
	sameColor := card.Color == top.Color
	sameType := card.Type == top.Type && card.Type != Number
	sameValue := card.Type == Number && top.Type == Number && card.Value == top.Value
	return sameColor || sameType || sameValue
}
