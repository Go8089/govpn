package udp

import (
	"fmt"
	"net"
)

type Client struct {
	ServerAddr string
}

func NewClient(addr string) *Client {
	return &Client{ServerAddr: addr}
}

func (c *Client) Send(message string) error {
	addr, err := net.ResolveUDPAddr("udp", c.ServerAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	buffer := make([]byte, 1024)

        n, _, err := conn.ReadFromUDP(buffer)
        if err != nil {
        return err}

        fmt.Println("Received:", string(buffer[:n]))
	return nil
}
