package udp

const (
	PacketPing uint8 = 1
	PacketPong uint8 = 2
)

type Packet struct {
	Version uint8
	Type    uint8
	Length  uint16
	Payload []byte
}

func (p *Packet) Marshal() []byte {
	data := make([]byte, 4+len(p.Payload))

	data[0] = p.Version
	data[1] = p.Type
	data[2] = byte(p.Length >> 8)
	data[3] = byte(p.Length)

	copy(data[4:], p.Payload)

	return data
}

func Unmarshal(data []byte) *Packet {
	return &Packet{
		Version: data[0],
		Type:    data[1],
		Length:  uint16(data[2])<<8 | uint16(data[3]),
		Payload: data[4:],
	}
}
