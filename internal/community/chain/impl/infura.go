package impl

import (
	"another_node/internal/community/chain"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

type InfuraCaller struct {
	rpc string
}

var _ chain.Caller = (*InfuraCaller)(nil)

type InfuraNetworks string

const (
	InfuraMainnet InfuraNetworks = "https://mainnet.infura.io/v3"
	InfuraSepolia InfuraNetworks = "https://sepolia.infura.io/v3"
)

func NewInfuraCaller(network InfuraNetworks, apiKey string) (chain.Caller, error) {
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("API Key is required")
	}

	return &InfuraCaller{
		rpc: fmt.Sprintf("%s/%s", string(network), apiKey),
	}, nil
}

func (a *InfuraCaller) connect() (*ethclient.Client, error) {
	conn, err := ethclient.Dial(a.rpc)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (a *InfuraCaller) Call(f func(*ethclient.Client) interface{}) error {
	if conn, err := a.connect(); err != nil {
		return err
	} else {
		defer conn.Close()
		ret := f(conn)
		if ret != nil {
			v, ok := ret.(error)
			if ok {
				return v
			}
		}
	}
	return nil
}
