package seedworks

type TxSignature struct {
	Nonce  string `json:"nonce" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Origin string `json:"origin" binding:"required"`
	TxData string `json:"txdata" binding:"required"`
}
