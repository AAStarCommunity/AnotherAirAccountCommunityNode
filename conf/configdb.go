package conf

import (
	"fmt"
	"sync"
)

var dsnTemplate = "host=%s port=%v user=%s password=%s dbname=%s TimeZone=%s sslmode=%s"

type ConfigDb struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	TimeZone string `yaml:"tz"`
	SslMode  string `yaml:"ssl_mode"`
}

var configDb *ConfigDb
var onceConfigDb sync.Once

func GetConfigDb() *ConfigDb {
	onceConfigDb.Do(func() {
		if configDb == nil {
			j := getConf().ConfigDb
			configDb = &ConfigDb{
				Host:     j.Host,
				Port:     j.Port,
				User:     j.User,
				Password: j.Password,
				DBName:   j.DBName,
				TimeZone: j.TimeZone,
				SslMode:  j.SslMode,
			}
		}
	})

	return configDb
}
func GetConfigDbDSN() string {
	confdb := GetConfigDb()
	return fmt.Sprintf(dsnTemplate, confdb.Host, confdb.Port, confdb.User, confdb.Password, confdb.DBName, confdb.TimeZone, confdb.SslMode)
}
