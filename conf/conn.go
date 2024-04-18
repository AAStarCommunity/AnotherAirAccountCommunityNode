package conf

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	User     string
	Password string
	Host     string
	Schema   string
}

var db *gorm.DB

func getDbConnectionString(c *Conf) *string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.Db.Host, c.Db.User, c.Db.Password, c.Db.Schema)
	return &dsn
}

var onceDb sync.Once

// GetDbClient 获取数据库连接对象
func GetDbClient() *gorm.DB {
	onceDb.Do(func() {
		if os.Getenv("UnitTestEnv") == "1" {
			db, _ = getInMemoryDbClient()
			db = db.Debug()
		} else {
			dsn := getDbConnectionString(getConf())
			_db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{
				DisableForeignKeyConstraintWhenMigrating: true,
			})
			if err != nil {
				panic(err)
			}

			if Environment.IsDevelopment() {
				_db = _db.Debug()
			}
			db = _db
		}
	})
	return db
}

// getInMemoryDbClient 获取内存数据库对象，仅限单元测试使用
func getInMemoryDbClient() (*gorm.DB, error) {
	if client, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
