package migrations

import (
	"another_node/conf"

	"gorm.io/gorm"
)

type Migration interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

var migrations []Migration

func init() {
	// Migration objects must be added in order from old to new
	migrations = []Migration{
		&Migration20240420{},
	}
}

func AutoMigrate() {
	db := conf.GetDbClient()

	// TODO: Check if there is a latest change record in the '__migration' table, if so, skip the migration
	migrate(db)
}

// migrate synchronizes database changes
func migrate(db *gorm.DB) {

	for i := 0; i < len(migrations); i++ {
		migrations[i].Up(db)
	}
}

// Rollback rolls back database changes
func Rollback(db *gorm.DB) {

	for i := len(migrations) - 1; i >= 0; i-- {
		migrations[i].Down(db)
	}
}
