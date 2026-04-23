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
	mutex   sync.RWMutex
	players map[string]*Player
}

var hub = &Hub{
	players: make(map[string]*Player),
}

func AddNewPlayer(nickname string, conn *websocket.Conn) {
	player := Player{
		Nickname: nickname,
		Conn:     conn,
	}

	hub.mutex.Lock()
	hub.players[player.Nickname] = &player
	hub.mutex.Unlock()
}

func GetPlayers() []string {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	nicknames := make([]string, 0, len(hub.players))

	for _, player := range hub.players {
		nicknames = append(nicknames, player.Nickname)
	}

	return nicknames
}

func GetPlayerConnections() map[string]*websocket.Conn {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	connections := make(map[string]*websocket.Conn, len(hub.players))

	for nickname, player := range hub.players {
		connections[nickname] = player.Conn
	}

	return connections
}
