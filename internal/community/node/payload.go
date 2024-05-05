package node

type Payload struct {
	Account    string
	PublicKey  string
	RpcAddress string
	RpcPort    uint16
	Version    uint32
}
