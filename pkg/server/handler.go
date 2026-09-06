package server

import (
	"log"
	"net"
	"sync/atomic"

	"hare/pkg/transport"
)

func handleMessage(s *Server, conn net.Conn, m *transport.Message) {
	switch m.Cmd {
	case transport.CmdPublish:
		s.engine.Publish(*m)
		atomic.AddUint64(&s.pubCount, 1)
	case transport.CmdSubscribe:
		ch := s.engine.Subscribe(string(m.Topic))
		go func() {
			defer s.engine.Unsubscribe(string(m.Topic), ch)

			atomic.AddInt64(&s.activeConsumers, 1)
			defer atomic.AddInt64(&s.activeConsumers, -1)

			for {
				msg, ok := <-ch
				if !ok {
					return
				}
				if _, err := conn.Write(msg); err != nil {
					return
				}
			}
		}()

	default:
		log.Printf("Unknown command: %d\n", m.Cmd)
	}
}

func handleConnection(s *Server, conn net.Conn) {
	defer conn.Close()
	for {
		m, err := transport.DecodeMessage(conn)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("received command: %s to topic %s, body size = %d byte(s)\n", m.Cmd, m.Topic, len(m.Body))
		handleMessage(s, conn, m)
	}
}
