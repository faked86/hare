package main

import (
	"fmt"
	"hare/pkg/transport"
	"log"
	"net"
)

const (
	port    int    = 8585
	tName   string = "metrics"
	bufSize int    = 256
)

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	topic := []byte(tName)
	tLen := len([]byte(topic))

	body := make([]byte, 0)
	pLen := len(body) + len(topic) + transport.TLenSize

	msg := transport.EncodeMessage(transport.Message{
		Cmd:   transport.CmdSubscribe,
		PLen:  uint32(pLen),
		TLen:  uint16(tLen),
		Topic: topic,
		Body:  body,
	})
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg, err := transport.DecodeMessage(conn)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("received msg: ", string(msg.Body))
	}
}
