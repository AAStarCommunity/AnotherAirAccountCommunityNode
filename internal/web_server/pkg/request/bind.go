package request

type Bind struct {
	Account    string `json:"account"`
	PublicKey  string `json:"publicKey"`
	RpcAddress string `json:"rpcAddress"`
	Version    int    `json:"version"`
}
