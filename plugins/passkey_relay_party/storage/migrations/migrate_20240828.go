package migrations

import (
	"another_node/plugins/passkey_relay_party/storage/model"

	"gorm.io/gorm"
)

type Migration20240828 struct {
}

var _ Migration = (*Migration20240828)(nil)

func (m *Migration20240828) Up(db *gorm.DB) error {

	if !db.Migrator().HasTable(&model.AirAccount{}) {
		if err := db.AutoMigrate(&model.AirAccount{}); err != nil {
			return err
		}
	}

	if !db.Migrator().HasTable(&model.AirAccountChain{}) {
		if err := db.AutoMigrate(&model.AirAccountChain{}); err != nil {
			return err
		}
	}

	if !db.Migrator().HasTable(&model.Passkey{}) {
		if err := db.AutoMigrate(&model.Passkey{}); err != nil {
			return err
		}
	}

	if !db.Migrator().HasTable(&model.HdWallet{}) {
		if err := db.AutoMigrate(&model.HdWallet{}); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migration20240828) Down(db *gorm.DB) error {
	if err := db.Migrator().DropTable(
		&model.HdWallet{},
	); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(
		&model.Passkey{},
	); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(
		&model.AirAccountChain{},
	); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(
		&model.AirAccount{},
	); err != nil {
		return err
	}

	return nil
}
