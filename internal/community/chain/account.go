package chain

import (
	"another_node/conf"
	"another_node/internal/community/account"
	"another_node/internal/seedworks"

	"github.com/pavankpdev/goaa"
	"golang.org/x/xerrors"
)

const salt int64 = 1

func CreateSmartAccount(wallet *account.HdWallet, network seedworks.Chain) (string, error) {
	pk := "0x" + wallet.PrivateKey()
	networkConfig := conf.GetNetworkConfigByNetwork(network)
	if networkConfig == nil {
		return "", xerrors.Errorf("Failed to get network config for network: %s", network)
	}
	entrypointAddress := networkConfig.V06EntryPointAddress
	factoryAddress := networkConfig.V06FactoryAddress
	rpcUrl := networkConfig.RpcUrl

	params := goaa.SmartAccountProviderParams{
		OwnerPrivateKey:            pk,
		RPC:                        rpcUrl,
		EntryPointAddress:          entrypointAddress,
		SmartAccountFactoryAddress: factoryAddress,
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
