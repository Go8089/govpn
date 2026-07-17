package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

type AES struct {
	gcm cipher.AEAD
}

func NewAES(key []byte) (*AES, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AES{
		gcm: gcm,
	}, nil
}

func (a *AES) Encrypt(plaintext, nonce []byte) ([]byte, error) {

	ciphertext := a.gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nil
}

func (a *AES) Decrypt(ciphertext, nonce []byte) ([]byte, error) {

	plaintext, err := a.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
