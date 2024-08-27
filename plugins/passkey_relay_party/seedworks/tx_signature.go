package seedworks

type TxSignature struct {
	Ticket string `json:"ticket"`
	Email  string `json:"-"`
	Origin string `json:"origin"`
	TxData string `json:"txdata"`
}

type TxSignatureResult struct {
	Code       int    `json:"code"`
	TxData     string `json:"txdata"`
	Sign       string `json:"sign"`
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
}
