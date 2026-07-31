package udp

import (
	"crypto/ecdh"
	"fmt"
	"net"
	"time"

	"github.com/Go8089/govpn/internal/crypto"
)

const readTimeout = 5 * time.Second

type Client struct {
	ServerAddr string
}

func NewClient(addr string) *Client {
	return &Client{
		ServerAddr: addr,
	}
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

	if err := conn.SetDeadline(time.Now().Add(readTimeout)); err != nil {
		return err
	}

	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		return err
	}

	if err := c.sendHello(conn); err != nil {
		return err
	}

	serverPublic, err := c.exchangeKeys(conn, keyPair)
	if err != nil {
		return err
	}

	sharedSecret, err := crypto.ComputeSharedSecret(
		keyPair.Private,
		serverPublic,
	)
	if err != nil {
		return err
	}

	aesCipher, err := crypto.NewAES(sharedSecret)
	if err != nil {
		return err
	}

	fmt.Println("Secure session established.")

	if err := c.sendEncryptedPing(conn, aesCipher, message); err != nil {
		return err
	}

	return nil
}

func (c *Client) sendHello(conn *net.UDPConn) error {
	packet := &Packet{
		Version: ProtocolVersion,
		Type:    PacketHello,
		Nonce:   make([]byte, NonceSize),
		Payload: []byte("HELLO"),
	}

	packet.Length = uint16(len(packet.Payload))

	if _, err := conn.Write(packet.Marshal()); err != nil {
		return err
	}

	reply, err := readPacket(conn)
	if err != nil {
		return err
	}

	if reply.Type != PacketHelloAck {
		return fmt.Errorf("expected HELLO_ACK")
	}

	fmt.Println("HELLO_ACK received")

	return nil
}

func (c *Client) exchangeKeys(
	conn *net.UDPConn,
	keyPair *crypto.KeyPair,
) (*ecdh.PublicKey, error) {

	packet := &Packet{
		Version: ProtocolVersion,
		Type:    PacketClientKey,
		Nonce:   make([]byte, NonceSize),
		Payload: keyPair.Public.Bytes(),
	}

	packet.Length = uint16(len(packet.Payload))

	if _, err := conn.Write(packet.Marshal()); err != nil {
		return nil, err
	}

	reply, err := readPacket(conn)
	if err != nil {
		return nil, err
	}

	if reply.Type != PacketServerKey {
		return nil, fmt.Errorf("expected SERVER_KEY")
	}

	return ecdh.X25519().NewPublicKey(reply.Payload)
}

func (c *Client) sendEncryptedPing(
	conn *net.UDPConn,
	aesCipher *crypto.AES,
	message string,
) error {

	nonce, err := crypto.GenerateNonce()
	if err != nil {
		return err
	}

	ciphertext, err := aesCipher.Encrypt([]byte(message), nonce)
	if err != nil {
		return err
	}

	packet := &Packet{
		Version: ProtocolVersion,
		Type:    PacketPing,
		Nonce:   nonce,
		Payload: ciphertext,
	}

	packet.Length = uint16(len(ciphertext))

	if _, err := conn.Write(packet.Marshal()); err != nil {
		return err
	}

	reply, err := readPacket(conn)
	if err != nil {
		return err
	}

	if reply.Type != PacketPong {
		return fmt.Errorf("expected PONG")
	}

	plaintext, err := aesCipher.Decrypt(reply.Payload, reply.Nonce)
	if err != nil {
		return err
	}

	fmt.Println("Received:", string(plaintext))

	return nil
}

func readPacket(conn *net.UDPConn) (*Packet, error) {
	buffer := make([]byte, 2048)

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	return Unmarshal(buffer[:n])
}