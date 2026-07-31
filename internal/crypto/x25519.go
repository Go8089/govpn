package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
)

var (
	ErrNilPrivateKey = errors.New("nil private key")
	ErrNilPublicKey  = errors.New("nil public key")
)

type KeyPair struct {
	Private *ecdh.PrivateKey
	Public  *ecdh.PublicKey
}

func GenerateKeyPair() (*KeyPair, error) {
	curve := ecdh.X25519()

	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Private: privateKey,
		Public:  privateKey.PublicKey(),
	}, nil
}

func ComputeSharedSecret(
	privateKey *ecdh.PrivateKey,
	publicKey *ecdh.PublicKey,
) ([]byte, error) {

	if privateKey == nil {
		return nil, ErrNilPrivateKey
	}

	if publicKey == nil {
		return nil, ErrNilPublicKey
	}

	return privateKey.ECDH(publicKey)
}

func MarshalPublicKey(publicKey *ecdh.PublicKey) ([]byte, error) {
	if publicKey == nil {
		return nil, ErrNilPublicKey
	}

	return publicKey.Bytes(), nil
}

func ParsePublicKey(data []byte) (*ecdh.PublicKey, error) {
	return ecdh.X25519().NewPublicKey(data)
}