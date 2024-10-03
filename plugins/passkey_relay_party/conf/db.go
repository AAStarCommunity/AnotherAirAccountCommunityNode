package conf

import (
	"another_node/conf"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var onceDb sync.Once
var dbClient *gorm.DB

// GetDbClient 获取数据库连接对象
func GetDbClient() *gorm.DB {
	onceDb.Do(func() {
		if os.Getenv("UnitTest") == "1" {
			dbClient, _ = getInMemoryDbClient()
			dbClient = dbClient.Debug()
		} else {
			_db, err := gorm.Open(postgres.Open(Get().DbConnection), &gorm.Config{
				DisableForeignKeyConstraintWhenMigrating: true,
			})
			if err != nil {
				panic(err)
			}

			if conf.Environment.IsDevelopment() {
				_db = _db.Debug()
			}
			dbClient = _db
		}
	})
	return dbClient
}

// getInMemoryDbClient for UnitTests only
func getInMemoryDbClient() (*gorm.DB, error) {
	if client, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
