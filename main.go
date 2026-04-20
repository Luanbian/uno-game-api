package main

import (
	"log"

	"github.com/Luanbian/uno-game-api/internal/server"
)

func main() {
	s := server.New(":8080")
	log.Fatal(s.Start())
}
