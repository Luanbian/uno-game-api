// Package server provides a server HTTP and register routes
package server

import (
	"log"
	"net/http"

	"github.com/Luanbian/uno-game-api/internal/health"
)

type Server struct {
	mux  *http.ServeMux
	addr string
}

func StartOn(addr string) {
	server := &Server{
		mux:  http.NewServeMux(),
		addr: addr,
	}

	server.register(health.Route, health.Handler)

	log.Fatal(server.start())
}

func (server *Server) register(pattern string, handler http.HandlerFunc) {
	server.mux.HandleFunc(pattern, handler)
}

func (server *Server) start() error {
	log.Printf("Server running on port %s", server.addr)
	return http.ListenAndServe(server.addr, server.mux)
}
