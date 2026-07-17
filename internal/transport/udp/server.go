package udp

import (
	"fmt"
	"net"
	"github.com/Go8089/govpn/internal/crypto"
)

type Server struct {
	Port int
}

func NewServer(port int) *Server {
	return &Server{Port: port}
}

func (s *Server) Start() error {
	addr := &net.UDPAddr{
		Port: s.Port,
		IP:   net.IPv4zero,
	}

	conn, err := net.ListenUDP("udp", addr)
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
	fmt.Printf("GoVPN UDP Server listening on %d\n", s.Port)

	buffer := make([]byte, 1024)
for {
	n, clientAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println(err)
		continue
	}

	packet, err := Unmarshal(buffer[:n])
	if err != nil {
		fmt.Println("Invalid packet:", err)
		continue
	}

	plaintext, err := aesCipher.Decrypt(packet.Payload, packet.Nonce)
	if err != nil {
		fmt.Println(err)
		continue
	}

	fmt.Println("Received:", string(plaintext))

	nonce, err := crypto.GenerateNonce()
	if err != nil {
		fmt.Println(err)
		continue
	}

	ciphertext, err := aesCipher.Encrypt([]byte("PONG"), nonce)
	if err != nil {
		fmt.Println(err)
		continue
	}

	response := &Packet{
		Version: 1,
		Type:    PacketPong,
		Length:  uint16(len(ciphertext)),
		Nonce:   nonce,
		Payload: ciphertext,
	}

	_, err = conn.WriteToUDP(response.Marshal(), clientAddr)
	if err != nil {
		fmt.Println(err)
	}
}
	
}
