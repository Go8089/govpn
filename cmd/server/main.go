package main

import (
	"log"

	"github.com/Go8089/govpn/internal/transport/udp"
)

func main() {
	server := udp.NewServer(51820)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
