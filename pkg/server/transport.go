package server

import "fmt"

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

type Message struct {
	cmd     Command
	payload []byte
}
