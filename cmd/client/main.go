package main

import (
	"log"

	"github.com/Go8089/govpn/internal/transport/udp"
)

func main() {
	client := udp.NewClient("127.0.0.1:51820")

	if err := client.Send("PING"); err != nil {
		log.Fatal(err)
	}
}
