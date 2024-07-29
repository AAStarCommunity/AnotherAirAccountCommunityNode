package seedworks

import mapset "github.com/deckarep/golang-set/v2"

type Chain string

const (
	EthereumMainnet Chain = "ethereum-mainnet"
	EthereumSepolia Chain = "ethereum-sepolia"
	OptimismMainnet Chain = "optimism-mainnet"
	OptimismSepolia Chain = "optimism-sepolia"
	ArbitrumOne     Chain = "arbitrum-one"
	ArbitrumNova    Chain = "arbitrum-nova"
	ArbitrumSpeolia Chain = "arbitrum-sepolia"
	ScrollMainnet   Chain = "scroll-mainnet"
	ScrollSepolia   Chain = "scroll-sepolia"
	StarketMainnet  Chain = "starknet-mainnet"
	StarketSepolia  Chain = "starknet-sepolia"
	BaseMainnet     Chain = "base-mainnet"
	BaseSepolia     Chain = "base-sepolia"
)

var ethereumAdaptableNetWork = mapset.NewSet(EthereumMainnet, EthereumSepolia, ArbitrumOne, ArbitrumNova, ArbitrumSpeolia, OptimismMainnet, OptimismSepolia, ScrollMainnet, ScrollSepolia, BaseMainnet, BaseSepolia)

func IsEthereumAdaptableNetWork(network Chain) bool {
	return ethereumAdaptableNetWork.Contains(network)
}
