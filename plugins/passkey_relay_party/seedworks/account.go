package seedworks

type AccountInfo struct {
	InitCode   string `json:"init_code"`
	EOA        string `json:"eoa"`
	AA         string `json:"aa"`
	Email      string `json:"email"`
	PrivateKey string `json:"private_key"`
}
