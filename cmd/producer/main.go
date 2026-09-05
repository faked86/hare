package main

import (
	"fmt"
	"hare/pkg/transport"
	"log"
	"net"
)

const (
	port        int    = 8585
	tName       string = "metrics"
	testBodyStr string = "test test"
)

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	topic := []byte(tName)
	tLen := len([]byte(topic))

	body := []byte(testBodyStr)
	pLen := len(body) + len(topic) + transport.TLenSize

	msg := transport.EncodeMessage(transport.Message{
		Cmd:   transport.CmdPublish,
		PLen:  uint32(pLen),
		TLen:  uint16(tLen),
		Topic: topic,
		Body:  body,
	})
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
}
