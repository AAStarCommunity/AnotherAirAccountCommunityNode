package conf

import (
	"another_node/internal/global_const"
	"encoding/json"
	"os"
)

type NetWorkConfig struct {
	ChainId              string `json:"chain_id"`
	IsTest               bool   `json:"is_test"`
	V06EntryPointAddress string `json:"v06_entrypoint_address"`
	V06FactoryAddress    string `json:"v06_factory_address"`
	V07EntryPointAddress string `json:"v07_entrypoint_address"`
	V07FactoryAddress    string `json:"v07_factory_address"`
	RpcUrl               string `json:"rpc_url"`
}

var networkConfigMap map[global_const.Network]*NetWorkConfig

func InitNetworkConfig(configPath string) {
	if configPath == "" {
		panic("config path is empty")
	}
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	var mapValue map[string]NetWorkConfig
	//var mapValue map[string]any
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&mapValue)
	if err != nil {
		panic(err)
	}
	networkConfigMap = make(map[global_const.Network]*NetWorkConfig)
	for key, value := range mapValue {
		if value.RpcUrl == "" {
			panic(key + " rpc url is empty")
		}
		if value.V06EntryPointAddress == "" {
			panic(key + " v06 entry point address is empty")
		}
		if value.V06FactoryAddress == "" {
			panic(key + " v06 factory address is empty")
		}
		networkConfigMap[global_const.Network(key)] = &value
	}

}

func GetNetworkConfigByNetwork(network global_const.Network) *NetWorkConfig {
	return networkConfigMap[network]
}
