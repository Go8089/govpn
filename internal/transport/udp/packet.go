package udp
import "errors"

var (
	ErrInvalidVersion = errors.New("invalid protocol version")
	ErrInvalidLength  = errors.New("invalid packet length")
	ErrInvalidType    = errors.New("invalid packet type")
	ErrPacketTooShort = errors.New("packet too short")
)
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

func (p *Packet) Validate() error {
	if p.Version != 1 {
		return ErrInvalidVersion
	}

	if p.Type != PacketPing && p.Type != PacketPong {
		return ErrInvalidType
	}

	if p.Length != uint16(len(p.Payload)) {
		return ErrInvalidLength
	}

	return nil
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

func Unmarshal(data []byte) (*Packet, error) {
	if len(data) < 4 {
		return nil, ErrPacketTooShort
	}

	packet := &Packet{
		Version: data[0],
		Type:    data[1],
		Length:  uint16(data[2])<<8 | uint16(data[3]),
		Payload: data[4:],
	}

	if err := packet.Validate(); err != nil {
		return nil, err
	}

	return packet, nil
}


