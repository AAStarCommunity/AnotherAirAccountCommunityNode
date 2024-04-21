package migrations

import (
	"another_node/internal/community/storage"

	"gorm.io/gorm"
)

type Migration20240420 struct {
}

func (m *Migration20240420) Up(db *gorm.DB) error {
	if !db.Migrator().HasTable(&storage.Member{}) {
		if err := db.AutoMigrate(&storage.Member{}); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migration20240420) Down(db *gorm.DB) error {
	if err := db.Migrator().DropTable(
		&storage.Member{},
	); err != nil {
		return err
	}

	return nil
}
