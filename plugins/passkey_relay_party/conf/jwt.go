package conf

type JWT struct {
	Security string `yaml:"security"`
	Realm    string `yaml:"realm"`
	IdKey    string `yaml:"id_key"`
}

// GetJwtKey Get JWT configuration
func GetJwtKey() *JWT {
	return &Get().Jwt
}
