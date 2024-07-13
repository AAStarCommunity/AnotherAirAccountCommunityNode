package storage

import (
	"another_node/conf"
	passkey_conf "another_node/plugins/passkey_relay_party/conf"
	"another_node/plugins/passkey_relay_party/seedworks"
	"another_node/plugins/passkey_relay_party/storage/model"
	"time"

	"gorm.io/gorm"
)

type PgsqlStorage struct {
}

var _ Db = (*PgsqlStorage)(nil)

func NewPgsqlStorage() *PgsqlStorage {
	return &PgsqlStorage{}
}

func (db *PgsqlStorage) Save(user *seedworks.User) error {
	client := conf.GetDbClient()
	sqlDB, err := client.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	if data, err := user.Marshal(); err != nil {
		return err
	} else {
		if encrypted, err := seedworks.Encrypt([]byte(passkey_conf.Get().VaultSecret), data); err != nil {
			return err
		} else {
			return client.Model(&model.User{}).Create(&model.User{
				Email:   user.GetEmail(),
				Rawdata: encrypted,
			}).Error
		}
	}
}

func (db *PgsqlStorage) Find(email string) (*seedworks.User, error) {
	client := conf.GetDbClient()
	sqlDB, err := client.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	user := model.User{}
	if err := client.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	} else {
		if data, err := seedworks.Decrypt([]byte(passkey_conf.Get().VaultSecret), &user.Rawdata); err != nil {
			return nil, err
		} else {
			return seedworks.UnmarshalUser(&data)
		}
	}
}

func (db *PgsqlStorage) SaveChallenge(email, captcha string) error {
	client := conf.GetDbClient()

	sqlDB, err := client.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	return client.Model(&model.CaptchaChallenge{}).Create(&model.CaptchaChallenge{
		Type:   model.Email,
		Object: email,
		Code:   captcha,
	}).Error
}

func (db *PgsqlStorage) Challenge(email, captcha string) bool {
	client := conf.GetDbClient()

	sqlDB, err := client.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	success := false
	err = client.Transaction(func(tx *gorm.DB) error {
		challenge := model.CaptchaChallenge{}
		if err := tx.
			Where("object = ? AND code = ? AND type = ?", email, captcha, model.Email).
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
