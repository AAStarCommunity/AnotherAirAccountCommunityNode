package seedworks

type TxSignature struct {
	Nonce  string `json:"nonce"`
	Email  string `json:"-"`
	Origin string `json:"origin"`
	TxData string `json:"txdata"`
}

type TxSignatureResult struct {
	Code       int    `json:"code"`
	TxData     string `json:"txdata"`
	Sign       string `json:"sign"`
	PrivateKey string `json:"privateKey"`
}
