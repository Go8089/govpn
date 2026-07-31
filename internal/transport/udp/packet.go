package udp

import (
	"encoding/binary"
	"errors"
)

var (
	ErrInvalidVersion = errors.New("invalid protocol version")
	ErrInvalidLength  = errors.New("invalid packet length")
	ErrInvalidType    = errors.New("invalid packet type")
	ErrPacketTooShort = errors.New("packet too short")
)

const (
	PacketPing uint8 = iota + 1
	PacketPong
	PacketHello
	PacketHelloAck
	PacketClientKey
	PacketServerKey
)

const (
	ProtocolVersion = 1
	NonceSize       = 12
	HeaderSize      = 4
	PacketSize      = HeaderSize + NonceSize
)

type Packet struct {
	Version uint8
	Type    uint8
	Length  uint16
	Nonce   []byte
	Payload []byte
}

func (p *Packet) Validate() error {
	if p.Version != ProtocolVersion {
		return ErrInvalidVersion
	}

	switch p.Type {
	case PacketPing,
		PacketPong,
		PacketHello,
		PacketHelloAck,
		PacketClientKey,
		PacketServerKey:
	default:
		return ErrInvalidType
	}

	if len(p.Nonce) != NonceSize {
		return ErrPacketTooShort
	}

	if p.Length != uint16(len(p.Payload)) {
		return ErrInvalidLength
	}

	return nil
}

func (p *Packet) Marshal() []byte {
	data := make([]byte, HeaderSize+NonceSize+len(p.Payload))

	data[0] = p.Version
	data[1] = p.Type

	binary.BigEndian.PutUint16(data[2:4], uint16(len(p.Payload)))

	copy(data[4:16], p.Nonce)
	copy(data[16:], p.Payload)

	return data
}

func Unmarshal(data []byte) (*Packet, error) {
	if len(data) < PacketSize {
		return nil, ErrPacketTooShort
	}

	length := binary.BigEndian.Uint16(data[2:4])

	if len(data) != PacketSize+int(length) {
		return nil, ErrInvalidLength
	}

	packet := &Packet{
		Version: data[0],
		Type:    data[1],
		Length:  length,
		Nonce:   append([]byte(nil), data[4:16]...),
		Payload: append([]byte(nil), data[16:]...),
	}

	if err := packet.Validate(); err != nil {
		return nil, err
	}

	return packet, nil
}
