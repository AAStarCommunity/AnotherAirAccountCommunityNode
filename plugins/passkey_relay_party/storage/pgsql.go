package storage

import (
	"another_node/conf"
	"another_node/internal/community/account"
	passkey_conf "another_node/plugins/passkey_relay_party/conf"
	"another_node/plugins/passkey_relay_party/seedworks"
	"another_node/plugins/passkey_relay_party/storage/model"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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

// CreateAccount create user account with init hdwallets
// won't throw error if account already exists
func (db *PgsqlStorage) CreateAccount(email string, wallets []account.HdWallet) error {

	modelWallet := make([]model.HdWallet, len(wallets))
	for i := range wallets {
		if plain, err := json.Marshal(wallets[i]); err != nil {
			return err
		} else {
			if vault, err := seedworks.Encrypt(db.vaultSecret, plain); err != nil {
				return err
			} else {
				modelWallet[i] = model.HdWallet{
					WalletVault: vault,
				}
			}
		}
	}

	return db.client.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("email = ?", email).First(&model.AirAccount{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				airAccount := model.AirAccount{
					Email:    email,
					HdWallet: modelWallet,
				}

				if err := tx.Omit("updated_at").Create(&airAccount).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// TODO: network is a factor for judging the wallet used or not
func updateWalletUsed(w []model.HdWallet, usedWalletId []int64) (*model.HdWallet, error) {
	for i := range w {
		used := false
		for j := range usedWalletId {
			if w[i].ID == usedWalletId[j] {
				used = true
				break
			}
		}
		if !used {
			return &w[i], nil
		}
	}
	return nil, seedworks.ErrNoAvailableWallet{}
}

// SaveAccounts create or update user accounts
func (db *PgsqlStorage) SaveAccounts(user *seedworks.User) error {
	// support email only for now
	return db.client.Transaction(func(tx *gorm.DB) error {
		email, _, _ := user.GetAccounts()
		airAccount := model.AirAccount{}
		if err := db.client.Preload(clause.Associations).Where("email = ?", email).First(&airAccount).Error; err != nil {
			// SaveAccounts only update user account, so if user not exists, return error
			return err
		}

		changed := false

		for _, v := range user.GetChains() {
			flag := false
			for j := range airAccount.AirAccountChains {
				if airAccount.AirAccountChains[j].ChainName == string(v.Name) &&
					strings.EqualFold(airAccount.AirAccountChains[j].Alias, v.Alias) {
					flag = true
					break
				}
			}
			if flag {
				continue
			}
			if w, err := updateWalletUsed(airAccount.HdWallet, func() []int64 {
				ids := make([]int64, 0)
				for i := range airAccount.AirAccountChains {
					if airAccount.AirAccountChains[i].ChainName == string(v.Name) {
						ids = append(ids, airAccount.AirAccountChains[i].FromWalletID)
					}
				}
				return ids
			}()); err != nil {
				return err
			} else {
				if err := tx.Save(&w).Error; err != nil {
					return err
				}
				airAccount.AirAccountChains = append(airAccount.AirAccountChains, model.AirAccountChain{
					InitCode:     v.InitCode,
					AA_Address:   v.AA_Addr,
					ChainName:    string(v.Name),
					Alias:        v.Alias,
					FromWalletID: w.ID,
				})
			}
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

		if changed {
			if err := tx.Omit("created_at", "email").Save(&airAccount).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *PgsqlStorage) FindUser(userHandler string) (*seedworks.User, error) {
	airaccount := model.AirAccount{}
	if err := db.client.Preload(clause.Associations).Where("email = ?", userHandler).First(&airaccount).Error; err != nil {
		return nil, err
	} else {
		return seedworks.MappingUser(&airaccount, func(vault *string) (string, error) {
			return seedworks.Decrypt(db.vaultSecret, vault)
		})
	}
}

func (db *PgsqlStorage) FindUserByPasskey(userHandler, credId string) (*seedworks.User, error) {
	airaccount := model.AirAccount{}
	if err := db.client.Preload(clause.Associations).Where("email = ? AND credential_id = ?", userHandler, credId).First(&airaccount).Error; err != nil {
		return nil, err
	} else {
		return seedworks.MappingUser(&airaccount, func(vault *string) (string, error) {
			return seedworks.Decrypt(db.vaultSecret, vault)
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
