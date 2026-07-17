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

	packet := &Packet{
		Version: 1,
		Type:    PacketPing,
		Length:  uint16(len(message)),
		Payload: []byte(message),
	}

	_, err = conn.Write(packet.Marshal())
	if err != nil {
		return err
	}
	buffer := make([]byte, 1024)

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return err
	}

	fmt.Println("Received:", string(buffer[:n]))
	return nil
}
