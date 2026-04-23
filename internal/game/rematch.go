// Package game rematch controls the rematch of the game
package game

import "math/rand"

func Rematch(players []string) error {
	currentGame, err := GetCurrentGameState()
	if err != nil {
		return err
	}

	currentGame.mutex.Lock()
	defer currentGame.mutex.Unlock()

	currentGame.Deck.cards = []Card{}
	currentGame.Hands = make(map[string][]Card, len(currentGame.Players))
	currentGame.DiscardPile = []Card{}
	shufflePlayers(players)
	currentGame.CurrentPlayer = 0
	currentGame.LastPlayer = -1
	currentGame.SaidUno = make(map[string]bool, len(currentGame.Players))
	currentGame.Winner = ""
	currentGame.Direction = 1

	return nil
}

func shufflePlayers(players []string) {
	rand.Shuffle(len(players), func(i int, j int) {
		players[i], players[j] = players[j], players[i]
	})
}
