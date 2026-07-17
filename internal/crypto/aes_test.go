package crypto

import "testing"

func TestAES(t *testing.T) {

	key, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	aesCipher, err := NewAES(key)
	if err != nil {
		t.Fatal(err)
	}

	nonce, err := GenerateNonce()
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("Hello GoVPN")

	ciphertext, err := aesCipher.Encrypt(plaintext, nonce)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := aesCipher.Decrypt(ciphertext, nonce)
	if err != nil {
		t.Fatal(err)
	}

	if string(decrypted) != string(plaintext) {
		t.Fatal("encryption/decryption failed")
	}
}