package udp

import (
	"fmt"
	"net"
	"github.com/Go8089/govpn/internal/crypto"
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
        
        key, err := crypto.LoadKey("configs/dev.key")
if err != nil {
	return err
}

aesCipher, err := crypto.NewAES(key)
if err != nil {
	return err
}

nonce, err := crypto.GenerateNonce()
if err != nil {
	return err
}

ciphertext, err := aesCipher.Encrypt([]byte(message), nonce)
if err != nil {
	return err
}

	packet := &Packet{
	Version: 1,
	Type:    PacketPing,
	Length:  uint16(len(ciphertext)),
	Nonce:   nonce,
	Payload: ciphertext,
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

	response, err := Unmarshal(buffer[:n])
if err != nil {
	return err
}

plaintext, err := aesCipher.Decrypt(response.Payload, response.Nonce)
if err != nil {
	return err
}

fmt.Println("Received:", string(plaintext))
	return nil
}
