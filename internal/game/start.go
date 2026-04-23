// Package game start init the game
package game

func StartGame(players []string) (*GameState, error) {
	hands := make(map[string][]Card, len(players))
	saidUno := make(map[string]bool, len(players))

	for _, nickname := range players {
		hands[nickname] = make([]Card, 0, 7)
		saidUno[nickname] = false
	}

	deck := NewDeck()
	err := deck.distribute(hands)
	if err != nil {
		return nil, err
	}

	var firstCard Card
	for {
		firstCard, err = deck.pickUpCard()
		if err != nil {
			return nil, err
		}

		if firstCard.Type == Number {
			break
		}

		deck.cards = append(deck.cards, firstCard)
		deck.shuffle()
	}

	discardPile := []Card{firstCard}

	gameState := &GameState{
		Deck:          deck,
		Hands:         hands,
		DiscardPile:   discardPile,
		Players:       players,
		CurrentPlayer: 0,
		LastPlayer:    -1,
		SaidUno:       saidUno,
		Winner:        "",
		Direction:     1,
		StackedCards:  0,
	}

	return gameState, nil
}
