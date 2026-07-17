package crypto

import (
	"fmt"
	"os"
)

func LoadKey(path string) ([]byte, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(key) != KeySize {
		return nil, fmt.Errorf("invalid key size: got %d bytes, expected %d", len(key), KeySize)
	}

	return key, nil
}