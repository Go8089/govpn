package crypto

import "testing"

func TestSharedSecret(t *testing.T) {

	alice, _ := GenerateKeyPair()
	bob, _ := GenerateKeyPair()

	s1, _ := ComputeSharedSecret(
		alice.Private,
		bob.Public,
	)

	s2, _ := ComputeSharedSecret(
		bob.Private,
		alice.Public,
	)

	if string(s1) != string(s2) {
		t.Fatal("shared secrets differ")
	}
}
