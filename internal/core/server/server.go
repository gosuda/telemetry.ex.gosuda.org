package server

import (
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	mux *httprouter.Router
}

func NewServer() *Server {
	return &Server{
		mux: httprouter.New(),
	}
}

func (s *Server) Serve(ln net.Listener) error {
	return http.Serve(ln, s.mux)
}

func (s *Server) Shutdown() {

}
