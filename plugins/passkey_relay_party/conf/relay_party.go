package conf

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var once sync.Once

type RelayParty struct {
	Mail struct {
		Host     string
		Tls      bool
		Port     int
		User     string
		Password string
		Replier  string
	}
	Jwt          JWT    `yaml:"jwt"`
	DbConnection string `yaml:"db_connection"` // db connection string
	VaultSecret  string `yaml:"vault_secret"`  // encrypt & decrypt data into/from db
	DVT          DVT    `yaml:"dvt"`
}

var config *RelayParty

// Get read config from env take precedence over configfile
func Get() *RelayParty {
	once.Do(func() {
		if config == nil {
			filePath := getConfFilePath()
			config = getConfiguration(filePath)
		}
	})
	return config
}

// getConfiguration read config from file
func getConfiguration(filePath *string) *RelayParty {
	if file, err := os.ReadFile(*filePath); err != nil {
		panic("conf lost")
	} else {
		c := RelayParty{}
		err := yaml.Unmarshal(file, &c)
		if err != nil {
			panic("conf lost")
		}
		return &c
	}
}

func getConfFilePath() *string {
	osPtah := os.Getenv("passkey_conf_path")
	if osPtah != "" {
		return &osPtah
	}
	envName := "prod"
	if len(os.Getenv("Env")) > 0 {
		envName = os.Getenv("Env")
	}
	path := fmt.Sprintf("plugins/passkey_relay_party/conf/appsettings.%s.yaml", strings.ToLower(envName))
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		path = "plugins/passkey_relay_party/conf/appsettings.yaml"
	}
	return &path
}
