package udp

import (
	"crypto/ecdh"
	"fmt"
	"net"

	"github.com/Go8089/govpn/internal/crypto"
	"github.com/Go8089/govpn/internal/session"
)

type Server struct {
	Port     int
	Sessions *session.Manager
}

func NewServer(port int) *Server {
	return &Server{
		Port:     port,
		Sessions: session.NewManager(5 * 60 * 1e9), // 5 minutes
	}
}

func (s *Server) Start() error {
	addr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: s.Port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("GoVPN UDP Server listening on %d\n", s.Port)

	buffer := make([]byte, 4096)

	serverKeyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		return err
	}

	clientPublicKeys := make(map[string]*ecdh.PublicKey)
	ciphers := make(map[string]*crypto.AES)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		packet, err := Unmarshal(buffer[:n])
		if err != nil {
			fmt.Println("invalid packet:", err)
			continue
		}

		switch packet.Type {

		case PacketHello:
			reply := &Packet{
				Version: ProtocolVersion,
				Type:    PacketHelloAck,
				Nonce:   make([]byte, NonceSize),
				Payload: []byte("HELLO_ACK"),
			}
			reply.Length = uint16(len(reply.Payload))

			if _, err := conn.WriteToUDP(reply.Marshal(), clientAddr); err != nil {
				fmt.Println(err)
			}

		case PacketClientKey:
			publicKey, err := ecdh.X25519().NewPublicKey(packet.Payload)
			if err != nil {
				fmt.Println(err)
				continue
			}

			clientPublicKeys[clientAddr.String()] = publicKey

			sharedSecret, err := crypto.ComputeSharedSecret(
				serverKeyPair.Private,
				publicKey,
			)
			if err != nil {
				fmt.Println(err)
				continue
			}

			aesCipher, err := crypto.NewAES(sharedSecret)
			if err != nil {
				fmt.Println(err)
				continue
			}

			ciphers[clientAddr.String()] = aesCipher

			s.Sessions.Create(clientAddr, sharedSecret)

			reply := &Packet{
				Version: ProtocolVersion,
				Type:    PacketServerKey,
				Nonce:   make([]byte, NonceSize),
				Payload: serverKeyPair.Public.Bytes(),
			}
			reply.Length = uint16(len(reply.Payload))

			if _, err := conn.WriteToUDP(reply.Marshal(), clientAddr); err != nil {
				fmt.Println(err)
			}

		case PacketPing:
			aesCipher, ok := ciphers[clientAddr.String()]
			if !ok {
				fmt.Println("unknown client:", clientAddr)
				continue
			}

			plaintext, err := aesCipher.Decrypt(packet.Payload, packet.Nonce)
			if err != nil {
				fmt.Println(err)
				continue
			}

			s.Sessions.Touch(clientAddr)

			fmt.Printf("[%s] %s\n", clientAddr, string(plaintext))

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

			reply := &Packet{
				Version: ProtocolVersion,
				Type:    PacketPong,
				Nonce:   nonce,
				Payload: ciphertext,
			}
			reply.Length = uint16(len(reply.Payload))

			if _, err := conn.WriteToUDP(reply.Marshal(), clientAddr); err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("unsupported packet type:", packet.Type)
		}

		s.Sessions.Cleanup()
	}
}
