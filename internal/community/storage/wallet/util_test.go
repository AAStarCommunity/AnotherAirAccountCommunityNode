package wallet_storage

import "testing"

func TestAesCrypto(t *testing.T) {
	key := []byte("1234567890123456")
	plaintext := []byte("hello world")
	ciphertext, err := crypto(plaintext, key)
	if err != nil {
		t.Fatal(err)
	}

	decryptedtext, err := decrypt(ciphertext, key)
	if err != nil {
		t.Fatal(err)
	}

	if string(plaintext) != string(decryptedtext) {
		t.Fatal("decrypted text not equal to plaintext")
	}
}
