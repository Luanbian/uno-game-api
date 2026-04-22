// Package hub provides a lobby game logic, connection and players
package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Nickname string
	Conn     *websocket.Conn
}

type Hub struct {
	mutex   sync.Mutex
	Players map[string]*Player
}

var hub = &Hub{
	Players: make(map[string]*Player),
}

func AddNewPlayer(nickname string, conn *websocket.Conn) {
	player := Player{
		Nickname: nickname,
		Conn:     conn,
	}

	hub.mutex.Lock()
	hub.Players[player.Nickname] = &player
	hub.mutex.Unlock()
}
