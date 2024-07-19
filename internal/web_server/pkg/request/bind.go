package request

type Bind struct {
	Account   string `json:"account"`
	PublicKey string `json:"publicKey"`
}
type Sign struct {
	TxHash string `json:"txhash"` // 16 Hex

}
