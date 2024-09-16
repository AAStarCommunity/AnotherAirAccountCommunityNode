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

// SaveAccounts create or update user accounts
func (db *PgsqlStorage) SaveAccounts(user *seedworks.User, chain consts.Chain, alias string) error {
	if walletMarshal, err := user.WalletMarshal(); err != nil {
		return err
	} else {
		if walletVault, err := seedworks.Encrypt(db.vaultSecret, walletMarshal); err != nil {
			return err
		} else {
			initCode, aaAddr := user.GetChainAddresses(chain, alias)

			// support email only for now
			return db.client.Transaction(func(tx *gorm.DB) error {
				email, _, _ := user.GetAccounts()
				exists := true
				airAccount := model.AirAccount{}
				if err := db.client.Preload(clause.Associations).Where("email = ?", email).First(&airAccount).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						exists = false
						airAccount = model.AirAccount{
							Email:            email,
							Passkeys:         make([]model.Passkey, 0),
							AirAccountChains: make([]model.AirAccountChain, 0),
						}

						airAccount.HdWallet = model.HdWallet{
							WalletVault: walletVault,
						}
					} else {
						return err
					}
				}

				changed := false

				for i, v := range user.GetChains() {
					flag := false
					for j := range airAccount.AirAccountChains {
						if airAccount.AirAccountChains[j].ChainName == string(i) &&
							airAccount.AirAccountChains[j].Alias == v.Alias {
							flag = true
							break
						}
					}
					if flag {
						continue
					}
					airAccount.AirAccountChains = append(airAccount.AirAccountChains, model.AirAccountChain{
						InitCode:   *initCode,
						AA_Address: *aaAddr,
						ChainName:  string(i),
						Alias:      alias,
					})
					changed = true
				}

				for i := range user.WebAuthnCredentials() {
					cred := user.WebAuthnCredentials()[i]
					credId := base64.URLEncoding.EncodeToString(cred.ID)
					flag := false
					for j := range airAccount.Passkeys {
						if airAccount.Passkeys[j].CredentialId == credId {
							flag = true
							break
						}
					}
					if flag {
						continue
					}
					rawdata, _ := json.Marshal(cred)
					passkey := model.Passkey{
						CredentialId: credId,
						PublicKey:    base64.URLEncoding.EncodeToString(cred.PublicKey),
						Algorithm:    strconv.Itoa(int(cred.Attestation.PublicKeyAlgorithm)),
						Origin:       "-",
						Rawdata:      string(rawdata),
					}
					airAccount.Passkeys = append(airAccount.Passkeys, passkey)
					changed = true
				}

				if exists {
					if changed {
						if err := tx.Omit("created_at", "email").Save(&airAccount).Error; err != nil {
							return err
						}
					} else {
						return seedworks.ErrUserAlreadyExists{}
					}
				} else {
					if err := tx.Model(&model.AirAccount{}).Create(&airAccount).Error; err != nil {
						return err
					}
				}

				return nil
			})
		}
	}
}

func (db *PgsqlStorage) FindUser(userHandler string) (*seedworks.User, error) {
	airaccount := model.AirAccount{}
	if err := db.client.Preload(clause.Associations).Where("email = ?", userHandler).First(&airaccount).Error; err != nil {
		return nil, err
	} else {
		return seedworks.MappingUser(&airaccount, func() (string, error) {
			return seedworks.Decrypt(db.vaultSecret, &airaccount.HdWallet.WalletVault)
		})
	}
}

func (db *PgsqlStorage) FindUserByPasskey(userHandler, credId string) (*seedworks.User, error) {
	airaccount := model.AirAccount{}
	if err := db.client.Preload(clause.Associations).Where("email = ? AND credential_id = ?", userHandler, credId).First(&airaccount).Error; err != nil {
		return nil, err
	} else {
		return seedworks.MappingUser(&airaccount, func() (string, error) {
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
