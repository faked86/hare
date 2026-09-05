package transport

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type Command byte

func (c Command) String() string {
	switch c {
	case CmdPublish:
		return "PUBLISH"
	case CmdSubscribe:
		return "SUBSCRIBE"
	default:
		return fmt.Sprintf("UNKNOWN_COMMAND(%d)", c)
	}
}

const (
	CmdPublish   Command = 1
	CmdSubscribe Command = 2
)

const (
	PLenSize = 4
	TLenSize = 2
)

type Message struct {
	Cmd   Command
	PLen  uint32
	TLen  uint16
	Topic []byte
	Body  []byte
}

func DecodeMessage(conn net.Conn) (*Message, error) {
	header := make([]byte, 5)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	cmd := header[0]
	length := binary.BigEndian.Uint32(header[1:5])

	var payload []byte

	if length < 2 {
		return nil, errors.New("packet payload too small")
	} else {
		payload = make([]byte, length)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			return nil, err
		}
	}

	topicLen := binary.BigEndian.Uint16(payload[0:2])
	if length < 2+uint32(topicLen) {
		return nil, errors.New("packet truncated: topic length exceeds payload")
	}
	topic := payload[2 : 2+topicLen]
	body := payload[2+topicLen:]

	return &Message{
		Cmd:   Command(cmd),
		PLen:  length,
		TLen:  topicLen,
		Topic: topic,
		Body:  body,
	}, nil
}

func EncodeMessage(m Message) []byte {
	buf := bytes.NewBuffer([]byte{byte(m.Cmd)})

	pLenB := make([]byte, PLenSize)
	binary.BigEndian.PutUint32(pLenB, m.PLen)
	buf.Write(pLenB)

	tLenB := make([]byte, TLenSize)
	binary.BigEndian.PutUint16(tLenB, m.TLen)
	buf.Write(tLenB)

	buf.Write(m.Topic)
	buf.Write(m.Body)

	return buf.Bytes()
}
