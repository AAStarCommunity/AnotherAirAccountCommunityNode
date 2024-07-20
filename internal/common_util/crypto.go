package common_util

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func EthereumSignHexStr(msg string, privateKey *ecdsa.PrivateKey) (string, error) {
	msgByte, err := DecodeStringWithPrefix(msg)
	if err != nil {
		return "", err
	}
	if hash, err := crypto.Sign(accounts.TextHash(msgByte), privateKey); err != nil {
		return "", err
	} else {
		return EncodeToHexStringWithPrefix(hash), nil
	}
}

func SignMessage(message string, privateKey *ecdsa.PrivateKey) (string, error) {
	messageHash := accounts.TextHash([]byte(message))

	signature, err := crypto.Sign(messageHash, privateKey)
	if err != nil {
		return "", err
	}

	signature[crypto.RecoveryIDOffset] += 27

	return hexutil.Encode(signature), nil
}

func DecodeStringWithPrefix(data string) ([]byte, error) {
	if data[:2] == "0x" {
		data = data[2:]
	}
	return hex.DecodeString(data)
}
func EncodeToHexStringWithPrefix(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}
