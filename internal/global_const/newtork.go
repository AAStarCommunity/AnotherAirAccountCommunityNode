package global_const

import mapset "github.com/deckarep/golang-set/v2"

type Network string

const (
	EthereumMainnet Network = "ethereum-mainnet"
	EthereumSepolia Network = "ethereum-sepolia"
	OptimismMainnet Network = "optimism-mainnet"
	OptimismSepolia Network = "optimism-sepolia"
	ArbitrumOne     Network = "arbitrum-one"
	ArbitrumNova    Network = "arbitrum-nova"
	ArbitrumSpeolia Network = "arbitrum-sepolia"
	ScrollMainnet   Network = "scroll-mainnet"
	ScrollSepolia   Network = "scroll-sepolia"
	StarketMainnet  Network = "starknet-mainnet"
	StarketSepolia  Network = "starknet-sepolia"
	BaseMainnet     Network = "base-mainnet"
	BaseSepolia     Network = "base-sepolia"
)

type NewWorkStack string

var ethereumAdaptableNetWork = mapset.NewSet(EthereumMainnet, EthereumSepolia, ArbitrumOne, ArbitrumNova, ArbitrumSpeolia, OptimismMainnet, OptimismSepolia, ScrollMainnet, ScrollSepolia, BaseMainnet, BaseSepolia)

func IsEthereumAdaptableNetWork(network Network) bool {
	return ethereumAdaptableNetWork.Contains(network)
}
