package crypto

import "crypto/rand"

const NonceSize = 12 // GCM nonce

func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, NonceSize)

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
