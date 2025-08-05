package server

import (
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	mux *httprouter.Router
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{
		mux: httprouter.New(),
	}
}

// RegisterHandler registers a handler for a specific method and path
func (s *Server) RegisterHandler(method, path string, handler httprouter.Handle) {
	s.mux.Handle(method, path, handler)
}

func (s *Server) Serve(ln net.Listener) error {
	return http.Serve(ln, s.mux)
}

func (s *Server) Shutdown() {
	// Shutdown logic can be implemented here
}
