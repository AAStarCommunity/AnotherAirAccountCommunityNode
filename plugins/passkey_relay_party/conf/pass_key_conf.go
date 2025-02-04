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
	Jwt          JWT    `yaml:"jwt"`
	DbConnection string `yaml:"db_connection"` // db connection string
	VaultSecret  string `yaml:"vault_secret"`  // encrypt & decrypt data into/from db
}

var config *PassKeyConf

// Get read config from env take precedence over configfile
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

			jwt__security := os.Getenv("jwt__security")
			jwt__realm := os.Getenv("jwt__realm")
			jwt__idkey := os.Getenv("jwt__idkey")

			filePath := getConfFilePath()
			confFile := getConfiguration(filePath)

			config = &PassKeyConf{
				DbConnection: func() string {
					if dbConnection == "" {
						return confFile.DbConnection
					}
					return dbConnection
				}(),
				VaultSecret: func() string {
					if len(vaultSecret) == 0 {
						return confFile.VaultSecret
					}
					return vaultSecret
				}(),
				Jwt: func() JWT {
					if jwt__security == "" {
						jwt__security = confFile.Jwt.Security
					}
					if jwt__realm == "" {
						jwt__realm = confFile.Jwt.Realm
					}
					if jwt__idkey == "" {
						jwt__idkey = confFile.Jwt.IdKey
					}
					return JWT{
						Security: jwt__security,
						Realm:    jwt__realm,
						IdKey:    jwt__idkey,
					}
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
