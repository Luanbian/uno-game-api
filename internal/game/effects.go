// Package game effects controls the effects of each card
package game

import "fmt"

func ApplyCardEffect(card Card, gs *GameState) error {
	if card.Type == Number {
		return nil
	}

	switch card.Type {
	case Jump:
		jumpNextPlayer(gs)
		return nil
	case Inverter:
		invertDirection(gs)
		return nil
	case Plustwo:
		return plusTwoNextPlayer(gs)
	case Plusfour:
		return plusFourNextPlayer(gs)
	case Joker:
		return nil
	default:
		return fmt.Errorf("card type %s not implemented yet", card.Type)
	}
}

func jumpNextPlayer(gs *GameState) {
	gs.CurrentPlayer = (gs.CurrentPlayer + gs.Direction + len(gs.Players)) % len(gs.Players)
}

func invertDirection(gs *GameState) {
	gs.Direction *= -1
}

func plusTwoNextPlayer(gs *GameState) error {
	nextPlayer := gs.Players[(gs.CurrentPlayer+gs.Direction+len(gs.Players))%len(gs.Players)]

	for range 2 {
		topCard, err := gs.Deck.pickUpCard()
		if err != nil {
			return err
		}

		err = addCardInHand(nextPlayer, topCard, gs)
		if err != nil {
			return err
		}
	}

	err := resetSaidUno(nextPlayer, gs)
	if err != nil {
		return err
	}

	jumpNextPlayer(gs)

	return nil
}

func plusFourNextPlayer(gs *GameState) error {
	nextPlayer := gs.Players[(gs.CurrentPlayer+gs.Direction+len(gs.Players))%len(gs.Players)]

	for range 4 {
		topCard, err := gs.Deck.pickUpCard()
		if err != nil {
			return err
		}

		err = addCardInHand(nextPlayer, topCard, gs)
		if err != nil {
			return err
		}
	}

	err := resetSaidUno(nextPlayer, gs)
	if err != nil {
		return err
	}

	jumpNextPlayer(gs)

	return nil
}

func SelectColor(nickname string, color Color) (*GameState, error) {
	validColors := map[Color]bool{Red: true, Green: true, Blue: true, Yellow: true}
	if !validColors[color] {
		return nil, fmt.Errorf("invalid color: %s", color)
	}

	currentGame, err := GetCurrentGameState()
	if err != nil {
		return nil, fmt.Errorf("selecting color: %w: ", err)
	}
	currentGame.mutex.Lock()
	defer currentGame.mutex.Unlock()

	if !isColorSelectPending(currentGame) {
		return nil, fmt.Errorf("no color selection pending")
	}
	if nickname != currentGame.Players[currentGame.LastPlayer] {
		return nil, fmt.Errorf("only player %s can select color", currentGame.Players[currentGame.LastPlayer])
	}

	currentGame.DiscardPile = append(currentGame.DiscardPile, Card{Color: color, Type: ColorSelect, Value: -1})
	nextTurn(currentGame)

	return currentGame, nil
}
