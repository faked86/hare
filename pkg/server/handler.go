package server

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

const DefaultTopicName = "default"

func readMessage(conn net.Conn) (*Message, error) {
	header := make([]byte, 5)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	cmd := header[0]
	length := binary.BigEndian.Uint32(header[1:5])

	var payload []byte

	if length > 0 {
		payload = make([]byte, length)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			return nil, err
		}
	} else {
		payload = nil
	}

	return &Message{
		cmd:     Command(cmd),
		payload: payload,
	}, nil
}

// func logMessage(m *Message) {
// 	var cmdStr string
// 	if m.cmd == CmdPublish {
// 		cmdStr = "PUBLISH"
// 	}
// 	if m.cmd == CmdSubscribe {
// 		cmdStr = "SUBSCRIBE"
// 	}
// 	log.Printf("received command: %s, payload size = %d byte(s)\n", cmdStr, m.pLen)
// }

func handleMessage(s *Server, conn net.Conn, m *Message) {
	switch m.cmd {
	case CmdPublish:
		s.engine.Publish(DefaultTopicName, m.payload)
	case CmdSubscribe:
		ch := s.engine.Subscribe(DefaultTopicName)
		go func() {
			defer s.engine.Unsubscribe(DefaultTopicName, ch)
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
		log.Printf("Unknown command: %d\n", m.cmd)
	}
}

func handleConnection(s *Server, conn net.Conn) {
	defer conn.Close()
	for {
		m, err := readMessage(conn)
		if err != nil {
			return
		}
		// logMessage(m)
		handleMessage(s, conn, m)
	}
}
