package node

import (
	"another_node/internal/community/storage"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/spaolacci/murmur3"
)

func getAddr() ([]byte, error) {
	if ins, err := storage.Open(); err != nil {
		return nil, err
	} else {
		defer ins.Close()

		db := ins.Instance

		if addr, err := db.Get([]byte("node:addr"), nil); err != nil {
			return nil, err
		} else {
			return addr, nil
		}
	}
}

// generateIdentity represents generating a public/private key pair for the identity of this node
func generateIdentity() ([]byte, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Encode private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "Community Node Private Key",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)

	// Encode public key to PEM format
	publicKey := &privateKey.PublicKey
	publicKeyPEM := &pem.Block{
		Type:  "Community Node Public Key",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}
	publicKeyBytes := pem.EncodeToMemory(publicKeyPEM)

	if ins, err := storage.Open(); err != nil {
		return nil, err
	} else {
		defer ins.Close()

		db := ins.Instance

		if err = db.Put([]byte("node:private_key"), privateKeyBytes, nil); err != nil {
			return nil, err
		}

		if err = db.Put([]byte("node:public_key"), publicKeyBytes, nil); err != nil {
			return nil, err
		}

		pub := string(publicKeyBytes)
		addr := extractIdenty(&pub)
		if err = db.Put([]byte("node:addr"), []byte(*addr), nil); err != nil {
			return nil, err
		} else {
			return []byte(*addr), nil
		}
	}
}

func extractIdenty(pub *string) *string {
	hash32 := murmur3.Sum64([]byte(*pub))
	rlt := fmt.Sprintf("0x%x", hash32)

	return &rlt
}
