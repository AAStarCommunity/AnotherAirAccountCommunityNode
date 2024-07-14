package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v2"
)

var once sync.Once

type Conf struct {
	DbConnection       string `yaml:"db_connection"`
	Web                Web
	Node               Node
	Storage            string
	Provider           Provider
	ChainNetworks      map[string]string `yaml:"chain_networks"`
	ApiKeyAccessEnable bool              `yaml:"api_key_access_enable"`
}

var conf *Conf

// getConf Read configuration
func getConf() *Conf {
	once.Do(func() {
		if conf == nil {
			filePath := getConfFilePath()
			conf = getConfiguration(filePath)
			fmt.Printf("AirAccount getConfPath: [%s]\r\n", *filePath)
		}
	})
	return conf
}
func IsApiKeyAccessEnable() bool {
	return getConf().ApiKeyAccessEnable
}

// getConfiguration represent get the config from env or file
// env will overwrite the file
func getConfiguration(filePath *string) *Conf {

	if file, err := os.ReadFile(*filePath); err != nil {
		return mappingEnvToConf(nil)
	} else {
		fmt.Println("getConfiguration" + *filePath)
		c := Conf{}
		_ = yaml.Unmarshal(file, &c)
		cjson, _ := json.Marshal(c)
		fmt.Println("getConfiguration JSON", string(cjson))
		return mappingEnvToConf(&c)

	}
}

func mappingEnvToConf(fileConf *Conf) (envConf *Conf) {
	envConf = &Conf{
		Web:                Web{},
		Node:               Node{},
		Provider:           Provider{},
		ApiKeyAccessEnable: true,
	}
	if fileConf != nil {
		envConf.DbConnection = fileConf.DbConnection
		envConf.ApiKeyAccessEnable = fileConf.ApiKeyAccessEnable
		envConf.ChainNetworks = fileConf.ChainNetworks
	}

	if storage := os.Getenv("storage"); len(storage) > 0 {
		envConf.Storage = storage
	} else if fileConf != nil {
		envConf.Storage = fileConf.Storage
	} else {
		envConf.Storage = "another.dat"
	}

	if web__port := os.Getenv("web__port"); len(web__port) > 0 {
		if port, err := strconv.Atoi(web__port); err == nil && port > 0 {
			envConf.Web.Port = port
		} else {
			panic("web__port is invalid")
		}
	} else if fileConf != nil {
		if fileConf.Web.Port > 0 {
			envConf.Web.Port = fileConf.Web.Port
		} else {
			panic("web.port is invalid")
		}
	}

	if node__standalone := os.Getenv("node__standalone"); len(node__standalone) > 0 {
		envConf.Node.Standalone = node__standalone == "true"
	} else if fileConf != nil {
		envConf.Node.Standalone = fileConf.Node.Standalone
	}
	if node__genesis := os.Getenv("node__genesis"); len(node__genesis) > 0 {
		envConf.Node.Genesis = node__genesis == "true"
	} else if fileConf != nil {
		envConf.Node.Genesis = fileConf.Node.Genesis
	}
	if node__externaladdr := os.Getenv("node__externaladdr"); len(node__externaladdr) > 0 {
		envConf.Node.ExternalAddr = node__externaladdr
	} else if fileConf != nil {
		envConf.Node.ExternalAddr = fileConf.Node.ExternalAddr
	}
	if node__externalport := os.Getenv("node__externalport"); len(node__externalport) > 0 {
		if port, err := strconv.ParseUint(node__externalport, 10, 16); err == nil && port > 0 {
			envConf.Node.ExternalPort = uint16(port)
		} else {
			panic("node__externalport is invalid")
		}
	} else if fileConf != nil {
		if fileConf.Node.ExternalPort > 0 {
			envConf.Node.ExternalPort = fileConf.Node.ExternalPort
		} else {
			panic("node.externalport is invalid")
		}
	}
	if node__bindaddr := os.Getenv("node__bindaddr"); len(node__bindaddr) > 0 {
		envConf.Node.BindAddr = node__bindaddr
	} else if fileConf != nil {
		envConf.Node.BindAddr = fileConf.Node.BindAddr
	}
	if node__bindport := os.Getenv("node__bindport"); len(node__bindport) > 0 {
		if port, err := strconv.ParseUint(node__bindport, 10, 16); err == nil && port > 0 {
			envConf.Node.BindPort = uint16(port)
		} else {
			panic("node__bindport is invalid")
		}
	} else if fileConf != nil {
		if fileConf.Node.BindPort > 0 {
			envConf.Node.BindPort = fileConf.Node.BindPort
		} else {
			panic("node.bindport is invalid")
		}
	}
	if provider_alchemy := os.Getenv("provider__alchemy"); len(provider_alchemy) > 0 {
		envConf.Provider.Alchemy = provider_alchemy
	} else if fileConf != nil {
		if len(fileConf.Provider.Alchemy) == 0 {
			panic("provider.alchemy is invalid")
		} else {
			envConf.Provider.Alchemy = fileConf.Provider.Alchemy
		}
	}

	return
}
