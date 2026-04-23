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

	firstCard, err := deck.pickUpCard()
	if err != nil {
		return nil, err
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
	}

	return gameState, nil
}
