package request

type Bind struct {
	Account   string `json:"account"`
	PublicKey string `json:"publicKey"`
}
