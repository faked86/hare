package main

import (
	"log"

	"hare/pkg/server"
)

const (
	port int = 8585
)

func main() {
	s, err := server.New(port)
	if err != nil {
		log.Fatal(err)
	}

	s.Listen()
}
