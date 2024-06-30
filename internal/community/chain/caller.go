package chain

import "github.com/ethereum/go-ethereum/ethclient"

type Caller interface {
	Call(f func(*ethclient.Client) interface{}) error
}
