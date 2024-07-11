package chain

import (
	"another_node/conf"
	"another_node/internal/community/account"
	"another_node/internal/global_const"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pavankpdev/goaa"
	"golang.org/x/xerrors"
	"math/big"
	"strings"
)

const salt int64 = 1
const creatAccountAbiJson = `
		[
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "owner",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "salt",
						"type": "uint256"
					}
				],
				"name": "createAccount",
				"outputs": [
					{
						"internalType": "contract SimpleAccount",
						"name": "ret",
						"type": "address"
					}
				],
				"stateMutability": "nonpayable",
				"type": "function"
			}
		]
	`

var creatAccountAbi abi.ABI

func init() {
	abiVar, err := abi.JSON(strings.NewReader(creatAccountAbiJson))
	if err != nil {
		panic(err)
	}
	creatAccountAbi = abiVar

}
func CreateSmartAccount(wallet *account.HdWallet, network global_const.Network) (accountAddress string, initCodeStr string, err error) {
	pk := "0x" + wallet.PrivateKey()
	networkConfig := conf.GetNetworkConfigByNetwork(network)
	if networkConfig == nil {
		return "", "", xerrors.Errorf("Failed to get network config for network: %s", network)
	}
	entrypointAddress := networkConfig.V06EntryPointAddress
	factoryAddressStr := networkConfig.V06FactoryAddress
	rpcUrl := networkConfig.RpcUrl

	params := goaa.SmartAccountProviderParams{
		OwnerPrivateKey:            pk,
		RPC:                        rpcUrl,
		EntryPointAddress:          entrypointAddress,
		SmartAccountFactoryAddress: factoryAddressStr,
	}

	client, err := goaa.NewSmartAccountProvider(params)

	if err != nil {
		return "", "", err
	}

	address, err := client.GetSmartAccountAddress(salt)
	if err != nil {
		return "", "", err
	}
	factoryAddress := common.HexToAddress(factoryAddressStr)
	initCodeByte, err := GetAccountInitCode(address, factoryAddress, salt)
	if err != nil {
		return "", "", err
	}
	initCodeStr = "Ox" + hex.EncodeToString(initCodeByte)
	return address.Hex(), initCodeStr, nil
}

func GetAccountInitCode(accountAddress common.Address, factoryAddress common.Address, salt int64) ([]byte, error) {
	data, err := creatAccountAbi.Pack("createAccount", accountAddress, big.NewInt(salt))
	if err != nil {
		return nil, xerrors.Errorf("error encoding function data: %v", err)
	}
	data = append(factoryAddress.Bytes(), data...)

	return data, nil
}
