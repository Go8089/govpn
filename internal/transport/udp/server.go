package udp

import (
	"fmt"
	"net"
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

	fmt.Printf("GoVPN UDP Server listening on %d\n", s.Port)

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		packet, err := Unmarshal(buffer[:n])
		if err != nil {
			fmt.Println("Invalid packet:", err)
			continue
		}

		fmt.Printf("Version=%d Type=%d Payload=%s\n",
			packet.Version,
			packet.Type,
			string(packet.Payload),
		)
		response := &Packet{
			Version: 1,
			Type:    PacketPong,
			Length:  4,
			Payload: []byte("PONG"),
		}
		_, err = conn.WriteToUDP(response.Marshal(), clientAddr)
		if err != nil {
			fmt.Println(err)
		}

	}
}
