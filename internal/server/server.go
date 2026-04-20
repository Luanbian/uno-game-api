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

func New(addr string) *Server {
	server := &Server{
		mux:  http.NewServeMux(),
		addr: addr,
	}

	server.Register(health.Route, health.Handler)

	return server
}

func (server *Server) Register(pattern string, handler http.HandlerFunc) {
	server.mux.HandleFunc(pattern, handler)
}

func (server *Server) Start() error {
	log.Printf("Server running on port %s", server.addr)
	return http.ListenAndServe(server.addr, server.mux)
}
