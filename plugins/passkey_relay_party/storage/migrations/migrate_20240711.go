package migrations

import (
	"another_node/plugins/passkey_relay_party/storage/model"

	"gorm.io/gorm"
)

type Migration20240711 struct {
}

var _ Migration = (*Migration20240711)(nil)

func (m *Migration20240711) Up(db *gorm.DB) error {

	if !db.Migrator().HasTable(&model.User{}) {
		if err := db.AutoMigrate(&model.User{}); err != nil {
			return err
		}
	}

	if !db.Migrator().HasTable(&model.CaptchaChallenge{}) {
		if err := db.AutoMigrate(&model.CaptchaChallenge{}); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migration20240711) Down(db *gorm.DB) error {
	if err := db.Migrator().DropTable(
		&model.User{},
	); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(
		&model.CaptchaChallenge{},
	); err != nil {
		return err
	}

	return nil
}
