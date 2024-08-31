package conf

import (
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbClient *gorm.DB
	dbOnce   = sync.Once{}
)

func GetDbClient() *gorm.DB {
	conn := getConf().DbConnection
	dbOnce.Do(func() {
		configDBVar, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		dbClient = configDBVar

		if Environment.IsDevelopment() {
			dbClient = dbClient.Debug()
		}
	})
	return dbClient
}
