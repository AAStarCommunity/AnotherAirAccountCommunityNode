package conf

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var once sync.Once

type PassKeyConf struct {
	Mail struct {
		Host     string
		Tls      bool
		Port     int
		User     string
		Password string
		Replier  string
	}
	RelayParty struct {
		DisplayName string `yaml:"display_name"`
		Id          string `yaml:"id"` // relay party id
	} `yaml:"relay_party"`
	DbConnection string `yaml:"db_connection"` // db connection string
	VaultSecret  string `yaml:"vault_secret"`  // encrypt & decrypt data into/from db
}

var config *PassKeyConf

// Get read config from env or file (env 1st)
func Get() *PassKeyConf {
	once.Do(func() {
		if config == nil {
			mailHost := os.Getenv("mail__host")
			mailTls := os.Getenv("mail__tls")
			mailPortStr := os.Getenv("mail__port")
			var mailPort int64 = 995
			var err error
			if mailPort, err = strconv.ParseInt(mailPortStr, 0, 0); err != nil {
				mailPort = 995
			}
			mailUser := os.Getenv("mail__user")
			mailPassword := os.Getenv("mail__password")
			replier := os.Getenv("mail__replier")

			dbConnection := os.Getenv("passkey_db_connection")
			vaultSecret := os.Getenv("passkey_vault_secret")

			filePath := getConfFilePath()
			confFile := getConfiguration(filePath)

			config = &PassKeyConf{
				RelayParty: struct {
					DisplayName string `yaml:"display_name"`
					Id          string `yaml:"id"` // relay party id
				}{
					DisplayName: confFile.RelayParty.DisplayName,
					Id:          confFile.RelayParty.Id,
				},
				DbConnection: func() string {
					if dbConnection == "" {
						return confFile.DbConnection
					}
					return dbConnection
				}(),
				VaultSecret: func() string {
					if vaultSecret == "" {
						return confFile.VaultSecret
					}
					return vaultSecret
				}(),
				Mail: struct {
					Host     string
					Tls      bool
					Port     int
					User     string
					Password string
					Replier  string
				}{
					Host: func() string {
						if mailHost == "" {
							return confFile.Mail.Host
						}
						return mailHost
					}(),
					Tls: func() bool {
						if mailTls == "" {
							return confFile.Mail.Tls
						}
						return strings.EqualFold("true", mailTls)
					}(),
					Port: func() int {
						if mailPortStr == "" {
							return confFile.Mail.Port
						}
						return int(mailPort)
					}(),
					User: func() string {
						if mailUser == "" {
							return confFile.Mail.User
						}
						return mailUser
					}(),
					Password: func() string {
						if mailPassword == "" {
							return confFile.Mail.Password
						}
						return mailPassword
					}(),
					Replier: func() string {
						if replier == "" {
							return confFile.Mail.Replier
						}
						return replier
					}(),
				},
			}
		}
	})
	return config
}

// getConfiguration read config from file
func getConfiguration(filePath *string) *PassKeyConf {
	if file, err := os.ReadFile(*filePath); err != nil {
		panic("conf lost")
	} else {
		c := PassKeyConf{}
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
