package conf

import (
	"sync"
)

type JWT struct {
	Security string `yaml:"security"`
	Realm    string `yaml:"realm"`
	IdKey    string `yaml:"id_key"`
}

var jwt *JWT

var onceJwt sync.Once

// GetJwtKey Get JWT configuration
func GetJwtKey() *JWT {
	onceJwt.Do(func() {
		if jwt == nil {
			jwt = &Get().Jwt
		}
	})

	return jwt
}
