package conf

import (
	"another_node/internal/seedworks"
	"encoding/json"
	"sync"
)

var onceChainNetwork sync.Once
var networkConfigMap map[seedworks.Chain]*ChainNetworkConfig

type ChainNetworkConfig struct {
	ChainId              string `json:"chain_id"`
	IsTest               bool   `json:"is_test"`
	V06EntryPointAddress string `json:"v06_entrypoint_address"`
	V06FactoryAddress    string `json:"v06_factory_address"`
	V07EntryPointAddress string `json:"v07_entrypoint_address"`
	V07FactoryAddress    string `json:"v07_factory_address"`
	RpcUrl               string `json:"rpc_url"`
}

func GetNetworkConfigByNetwork(network seedworks.Chain) *ChainNetworkConfig {
	onceChainNetwork.Do(func() {
		if len(networkConfigMap) == 0 {
			j := getConf().ChainNetworks
			networkConfigMap = make(map[seedworks.Chain]*ChainNetworkConfig)
			for key, value := range j {
				var v ChainNetworkConfig
				_ = json.Unmarshal([]byte(value), &v)
				if v.RpcUrl == "" {
					panic(key + " rpc url is empty")
				}
				if v.V06EntryPointAddress == "" {
					panic(key + " v06 entry point address is empty")
				}
				if v.V06FactoryAddress == "" {
					panic(key + " v06 factory address is empty")
				}
				networkConfigMap[seedworks.Chain(key)] = &v
			}
		}
	})

	return networkConfigMap[network]
}
