// Package game deck provides a basic deck and shuffle
package game

import (
	"fmt"
	"math/rand"
)

type Deck struct {
	cards []Card
}

func NewDeck() *Deck {
	deck := &Deck{
		cards: make([]Card, 0, 108),
	}

	colors := []Color{Blue, Green, Red, Yellow}
	basicTypes := []Type{Jump, Inverter, Plustwo}
	wildTypes := []Type{Joker, Plusfour}

	for i := 0; i <= 9; i++ {
		for _, color := range colors {
			c := Card{Color: color, Type: Number, Value: i}

			deck.cards = append(deck.cards, c)
			if i != 0 {
				deck.cards = append(deck.cards, c)
			}
		}
	}

	for _, color := range colors {
		for _, basicType := range basicTypes {
			c := Card{Color: color, Type: basicType, Value: -1}

			deck.cards = append(deck.cards, c, c)
		}
	}

	for _, wildType := range wildTypes {
		c := Card{Color: None, Type: wildType, Value: -1}

		deck.cards = append(deck.cards, c, c, c, c)
	}

	deck.shuffle()
	return deck
}

func (deck *Deck) shuffle() {
	rand.Shuffle(len(deck.cards), func(i int, j int) {
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
	})
}

func (deck *Deck) pickUpCard() (Card, error) {
	if len(deck.cards) == 0 {
		return Card{}, fmt.Errorf("deck is empty")
	}

	first := deck.cards[0]
	deck.cards = deck.cards[1:]
	return first, nil
}

func (deck *Deck) distribute(hands map[string][]Card) error {
	if len(hands) == 0 || len(deck.cards) == 0 {
		return fmt.Errorf("no players connected or deck is empty")
	}
	if len(deck.cards) < len(hands)*7 {
		return fmt.Errorf("not enough cards to distribute")
	}

	for nickname := range hands {
		currentHand := deck.cards[0:7]
		deck.cards = deck.cards[7:]
		hands[nickname] = append(hands[nickname], currentHand...)
	}

	return nil
}
