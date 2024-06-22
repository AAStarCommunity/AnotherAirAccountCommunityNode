package conf

import (
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v2"
)

var once sync.Once

type Conf struct {
	Web      Web
	Jwt      JWT
	Node     Node
	Storage  string
	Provider Provider
}

var conf *Conf

// getConf Read configuration
func getConf() *Conf {
	once.Do(func() {
		if conf == nil {
			filePath := getConfFilePath()
			conf = getConfiguration(filePath)
		}
	})
	return conf
}

// getConfiguration represent get the config from env or file
// env will overwrite the file
func getConfiguration(filePath *string) *Conf {
	if file, err := os.ReadFile(*filePath); err != nil {
		return mappingEnvToConf(nil)
	} else {
		c := Conf{}
		yaml.Unmarshal(file, &c)
		return mappingEnvToConf(&c)
	}
}

func mappingEnvToConf(fileConf *Conf) (envConf *Conf) {
	envConf = &Conf{
		Web:      Web{},
		Jwt:      JWT{},
		Node:     Node{},
		Provider: Provider{},
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

	if jwt__security := os.Getenv("jwt__security"); len(jwt__security) > 0 {
		envConf.Jwt.Security = jwt__security
	} else if fileConf != nil {
		envConf.Jwt.Security = fileConf.Jwt.Security
	}
	if jwt__realm := os.Getenv("jwt__realm"); len(jwt__realm) > 0 {
		envConf.Jwt.Security = jwt__realm
	} else if fileConf != nil {
		envConf.Jwt.Realm = fileConf.Jwt.Realm
	}
	if jwt__idkey := os.Getenv("jwt__idkey"); len(jwt__idkey) > 0 {
		envConf.Jwt.Security = jwt__idkey
	} else if fileConf != nil {
		envConf.Jwt.IdKey = fileConf.Jwt.IdKey
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
