package server

import (
	"fmt"
	"log"
	"net"

	"hare/pkg/server/engine"
)

type Server struct {
	port   int
	ln     net.Listener
	engine *engine.Engine
}

func New(port int) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	e := engine.New()
	return &Server{
		port:   port,
		ln:     ln,
		engine: e,
	}, nil
}

func (s *Server) Listen() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(s, conn)
	}
}
