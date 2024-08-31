package storage

import (
	"another_node/conf"
	consts "another_node/internal/seedworks"
	passkey_conf "another_node/plugins/passkey_relay_party/conf"
	"another_node/plugins/passkey_relay_party/seedworks"
	"another_node/plugins/passkey_relay_party/storage/model"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PgsqlStorage struct {
	client      *gorm.DB
	vaultSecret []byte
}

var _ Db = (*PgsqlStorage)(nil)

func NewPgsqlStorage() *PgsqlStorage {
	return &PgsqlStorage{
		client:      conf.GetDbClient(),
		vaultSecret: []byte(passkey_conf.Get().VaultSecret),
	}
}

func (db *PgsqlStorage) SaveAccounts(user *seedworks.User, chain consts.Chain) error {
	if walletMarshal, err := user.WalletMarshal(); err != nil {
		return err
	} else {
		if walletVault, err := seedworks.Encrypt(db.vaultSecret, walletMarshal); err != nil {
			return err
		} else {
			initCode, aaAddr, eoaAddr := user.GetChainAddresses(chain)

			// so far, support email only
			return db.client.Transaction(func(tx *gorm.DB) error {
				email, _, _ := user.GetAccounts()
				if exists, err := db.FindUser(email); err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}

					if exists != nil {
						return seedworks.ErrUserAlreadyExists{}
					}

					newAirAccount := model.AirAccount{
						Email:            email,
						Passkeys:         make([]model.Passkey, 0),
						AirAccountChains: make([]model.AirAccountChain, 0),
					}

					newAirAccount.HdWallet = model.HdWallet{
						WalletVault: walletVault,
					}

					for i := range user.WebAuthnCredentials() {
						cred := user.WebAuthnCredentials()[i]
						rawdata, _ := json.Marshal(cred)
						passkey := model.Passkey{
							CredentialId: base64.URLEncoding.EncodeToString(cred.ID),
							PublicKey:    base64.URLEncoding.EncodeToString(cred.PublicKey),
							Algorithm:    strconv.Itoa(int(cred.Attestation.PublicKeyAlgorithm)),
							Origin:       "-",
							Rawdata:      string(rawdata),
						}
						newAirAccount.Passkeys = append(newAirAccount.Passkeys, passkey)
					}

					newAirAccount.AirAccountChains = append(newAirAccount.AirAccountChains, model.AirAccountChain{
						InitCode:    *initCode,
						AA_Address:  *aaAddr,
						EOA_Address: *eoaAddr,
						ChainName:   string(chain),
					})

					if err := tx.Model(&model.AirAccount{}).Create(&newAirAccount).Error; err != nil {
						return err
					}
				}
				return nil
			})
		}
	}
}

func (db *PgsqlStorage) FindUser(email string) (*seedworks.User, error) {
	airaccount := model.AirAccount{}
	if err := db.client.Preload(clause.Associations).Where("email = ?", email).First(&airaccount).Error; err != nil {
		return nil, err
	} else {
		return seedworks.NewUser(&airaccount, func() (string, error) {
			return seedworks.Decrypt(db.vaultSecret, &airaccount.HdWallet.WalletVault)
		})
	}
}

func (db *PgsqlStorage) SaveChallenge(captchaType model.ChallengeType, email, captcha string) error {
	return db.client.Model(&model.CaptchaChallenge{}).Create(&model.CaptchaChallenge{
		Type:   captchaType,
		Object: email,
		Code:   captcha,
	}).Error
}

func (db *PgsqlStorage) Challenge(captchaType model.ChallengeType, email, captcha string) bool {
	success := false
	err := db.client.Transaction(func(tx *gorm.DB) error {
		challenge := model.CaptchaChallenge{}
		if err := tx.
			Where("object = ? AND code = ? AND type = ?", email, captcha, captchaType).
			Order("created_at DESC").
			First(&challenge).Error; err != nil {
			return err
		} else {
			if challenge.ID > 0 {
				success = time.Since(challenge.CreatedAt) < 10*time.Minute
				return tx.Model(&model.CaptchaChallenge{}).Where("id = ?", challenge.ID).Delete(&challenge).Error
			}
		}
		return nil
	})

	return err == nil && success
}

func (db *PgsqlStorage) GetAccountsByEmail(email, chain string) (initCode, addr, eoaAddr string, err error) {
	account := model.UserAccount{}
	if err := db.client.Where("email = ? AND chain = ?", email, chain).First(&account).Error; err != nil {
		return "", "", "", err
	} else {
		return account.InitCode, account.Address, account.EoaAddress, nil
	}
}
