package node

import (
	"another_node/conf"
	"another_node/internal/community/storage"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/spaolacci/murmur3"
	"github.com/syndtr/goleveldb/leveldb"
)

func getAddr() ([]byte, error) {
	if db, err := storage.EnsureOpen(); err != nil {
		return nil, err
	} else {
		if addr, err := db.Get([]byte("node:addr"), nil); err != nil {
			if errors.Is(err, leveldb.ErrNotFound) {
				return generateIdentity(db)
			} else {
				return nil, err
			}
		} else {
			return addr, nil
		}
	}
}

// generateIdentity represents generating a public/private key pair for the identity of this node
func generateIdentity(db *leveldb.DB) ([]byte, error) {
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
		nodeAddr := &storage.NodeAddr{
			Addr:     *addr,
			Endpoint: conf.GetNode().ExternalAddr + ":" + fmt.Sprint(conf.GetNode().ExternalPort),
		}
		db.Put([]byte(storage.NodeKey(nodeAddr)), nodeAddr.Marshal(), nil)

		return []byte(*addr), nil
	}
}

func extractIdenty(pub *string) *string {
	hash32 := murmur3.Sum64([]byte(*pub))
	rlt := fmt.Sprintf("0x%x", hash32)

	return &rlt
}
