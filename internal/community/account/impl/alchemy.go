package impl

import (
	"another_node/internal/community/account"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
	"github.com/pavankpdev/goaa"
)

type AlchemyProvider struct {
	rpc string
}

var _ account.Provider = (*AlchemyProvider)(nil)

const salt int64 = 1

const EntryPointContractAddress = "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"

// see ref: https://docs.alchemy.com/reference/factory-addresses#deployment-addresses
const AlchemyWalletContractFactory = "0x00004EC70002a32400f8ae005A26081065620D20"

func NewAlchemyProvider(apiKey string) (account.Provider, error) {
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("API Key is required")
	}
	return &AlchemyProvider{
		rpc: "https://eth-sepolia.g.alchemy.com/v2/" + apiKey,
	}, nil
}

func (a *AlchemyProvider) GetRpc() string {
	return a.rpc
}

func (a *AlchemyProvider) CreateAccount(wallet *account.HdWallet) (string, error) {

	pk := "0x" + wallet.PrivateKey()

	eth := big.NewFloat(0.01)
	wei := new(big.Float)
	wei.Mul(eth, big.NewFloat(params.Ether))

	params := goaa.SmartAccountProviderParams{
		OwnerPrivateKey:            pk,
		RPC:                        a.rpc,
		EntryPointAddress:          EntryPointContractAddress,
		SmartAccountFactoryAddress: AlchemyWalletContractFactory,
	}

	client, err := goaa.NewSmartAccountProvider(params)

	if err != nil {
		return "", err
	}

	address, err := client.GetSmartAccountAddress(salt)
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}
