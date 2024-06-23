package wallet_storage

import (
	"log"
)

// Wallet represents a raw data of wallet in storage(LevelDb)
type Wallet struct {
	mnemonic   string // wallet mnemonic
	Address    string // wallet address
	privateKey string // wallet private key
	AAAdress   string // aa address from wallet
}

const (
	mnemonicCap         = 512
	addressCap          = 128
	privateKeyCap       = 128
	aaAddressCap        = 128
	walletMarshalHeader = byte(0x01)
)

const secretKey = "aastar@Planckeraaaaaaaaa"
const walletPrefix = "wallet:"

// marshalTotalCap is the total cap of marshaled wallet
// 1 byte for header
// 512 bytes for mnemonic
// 128 bytes for address
// 128 bytes for private key
// 128 bytes for aa address
var marshalTotalCap = 1 + mnemonicCap + addressCap + privateKeyCap + aaAddressCap

// Marshal convert wallet to bytes with encryption
func (w *Wallet) Marshal() []byte {

	mnemonicBytes := []byte(w.mnemonic)
	if len(mnemonicBytes) > mnemonicCap {
		return nil
	}

	addressBytes := []byte(w.Address)
	if len(addressBytes) > addressCap {
		return nil
	}

	privateKeyBytes := []byte(w.privateKey)
	if len(privateKeyBytes) > privateKeyCap {
		return nil
	}

	aaAddressBytes := []byte(w.AAAdress)
	if len(aaAddressBytes) > aaAddressCap {
		return nil
	}

	ret := make([]byte, 0, marshalTotalCap)
	offset := 0
	copy(ret, []byte{walletMarshalHeader})
	offset += 1
	copy(ret[offset:offset+mnemonicCap], mnemonicBytes)
	offset += mnemonicCap
	copy(ret[offset:offset+addressCap], addressBytes)
	offset += addressCap
	copy(ret[offset:offset+privateKeyCap], privateKeyBytes)
	offset += privateKeyCap
	copy(ret[offset:offset+aaAddressCap], aaAddressBytes)

	if c, err := crypto(ret, []byte(secretKey)); err != nil {
		log.Default().Println("crypto error: ", err)
		return nil
	} else {
		return c
	}
}

func (w *Wallet) Unmarshal(data []byte) error {
	if len(data) < marshalTotalCap {
		return ErrInvalidWalletData{}
	}

	if d, err := decrypt(data, []byte(secretKey)); err != nil {
		log.Default().Println("decrypt error: ", err)
		return err
	} else {
		if d[0] != walletMarshalHeader {
			return ErrInvalidWalletData{}
		}

		offset := 1
		w.mnemonic = string(d[offset : offset+mnemonicCap])
		offset += mnemonicCap
		w.Address = string(d[offset : offset+addressCap])
		offset += addressCap
		w.privateKey = string(d[offset : offset+privateKeyCap])
		offset += privateKeyCap
		w.AAAdress = string(d[offset : offset+aaAddressCap])
	}

	return nil
}

func (w *Wallet) Key() string {
	return walletPrefix + w.Address
}
