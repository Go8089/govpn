package crypto

import "crypto/rand"

const KeySize = 32 // AES-256

func GenerateKey() ([]byte, error) {
	key := make([]byte, KeySize)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
