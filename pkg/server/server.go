package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"hare/pkg/server/engine"
)

type Server struct {
	port            int
	ln              net.Listener
	engine          *engine.Engine
	pubCount        uint64
	activeConsumers int64
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
	go s.startMetricsLogger()
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(s, conn)
	}
}

func (s *Server) startMetricsLogger() {
	for {
		log.Printf("Server is running. Producers: %d Consumers: %d", atomic.LoadUint64(&s.pubCount), atomic.LoadInt64(&s.activeConsumers))
		time.Sleep(5 * time.Second)
	}
}
