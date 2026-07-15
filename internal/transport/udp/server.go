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
		response := []byte("PONG")

                _, err = conn.WriteToUDP(response, clientAddr)
               if err != nil {
               fmt.Println(err)}

		fmt.Printf("Received %d bytes from %s: %s\n",
			n,
			clientAddr.String(),
			string(buffer[:n]),
		)
	}
}
