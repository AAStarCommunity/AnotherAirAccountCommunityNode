package seedworks

type TxSignature struct {
	Nonce  string `json:"nonce" binding:"required"`
	Email  string `json:"-"`
	Origin string `json:"origin" binding:"required"`
	TxData string `json:"txdata" binding:"required"`
}

type TxSignatureResult struct {
	Code   int    `json:"code"`
	TxData string `json:"txdata"`
	Sign   string `json:"sign"`
}
