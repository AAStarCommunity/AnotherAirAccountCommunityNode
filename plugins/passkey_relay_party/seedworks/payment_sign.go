package seedworks

type PaymentSign struct {
	Nonce  string `json:"nonce" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Origin string `json:"origin" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}
