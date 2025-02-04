package wallet_storage

import (
	"another_node/internal/community/storage"
	"log"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/util"
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

const secretKey = "aastar@Plancker^"
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

	ret := make([]byte, marshalTotalCap)
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

	s := string(ret)
	_ = s
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
		w.mnemonic = strings.Trim(string(d[offset:offset+mnemonicCap]), "\x00")
		offset += mnemonicCap
		w.Address = strings.Trim(string(d[offset:offset+addressCap]), "\x00")
		offset += addressCap
		w.privateKey = strings.Trim(string(d[offset:offset+privateKeyCap]), "\x00")
		offset += privateKeyCap
		w.AAAdress = strings.Trim(string(d[offset:offset+aaAddressCap]), "\x00")
	}

	return nil
}

func (w *Wallet) Key(hashedAccount *string) string {
	return walletPrefix + *hashedAccount + ":" + w.Address
}

func UpsertWallet(hashedAccount *string, wallet *Wallet) error {
	if wallet == nil {
		return ErrNilWallet{}
	}
	if db, err := storage.EnsureOpen(); err != nil {
		return err
	} else {
		return db.Put([]byte(wallet.Key(hashedAccount)), wallet.Marshal(), nil)
	}
}

func TryFindWallet(hashedAccount string) ([]Wallet, error) {
	if db, err := storage.EnsureOpen(); err != nil {
		return nil, err
	} else {
		iter := db.NewIterator(
			util.BytesPrefix([]byte(walletPrefix+hashedAccount+":")),
			nil)
		defer iter.Release()

		var wallets []Wallet
		for iter.Next() {
			wallet := Wallet{}
			if err := wallet.Unmarshal(iter.Value()); err != nil {
				return nil, err
			}
			wallets = append(wallets, wallet)
		}

		return wallets, nil
	}
}
