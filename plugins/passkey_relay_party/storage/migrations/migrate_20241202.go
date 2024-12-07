package migrations

import (
	"another_node/plugins/passkey_relay_party/storage/model"

	"gorm.io/gorm"
)

type Migration20241202 struct {
}

var _ Migration = (*Migration20241202)(nil)

func (m *Migration20241202) Up(db *gorm.DB) error {

	if !db.Migrator().HasColumn(&model.AirAccount{}, "eoa_address") {
		if err := db.Migrator().AddColumn(&model.AirAccount{}, "eoa_address"); err != nil {
			return err
		}
	}

	if !db.Migrator().HasColumn(&model.AirAccount{}, "webauthn_id") {
		if err := db.Migrator().AddColumn(&model.AirAccount{}, "webauthn_id"); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migration20241202) Down(db *gorm.DB) error {
	if db.Migrator().HasColumn(&model.AirAccount{}, "eoa_address") {
		if err := db.Migrator().DropColumn(&model.AirAccount{}, "eoa_address"); err != nil {
			return err
		}
	}

	if db.Migrator().HasColumn(&model.AirAccount{}, "webauthn_id") {
		if err := db.Migrator().DropColumn(&model.AirAccount{}, "webauthn_id"); err != nil {
			return err
		}
	}

	return nil
}
