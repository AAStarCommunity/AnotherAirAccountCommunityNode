package conf

import (
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v2"
)

var once sync.Once

type Conf struct {
	Web  Web
	Db   DB
	Jwt  JWT
	Node Node
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
		err := yaml.Unmarshal(file, &c)
		if err != nil {
			return mappingEnvToConf(&c)
		}

		return &c
	}
}

func mappingEnvToConf(fileConf *Conf) *Conf {
	envConf := &Conf{
		Web:  Web{},
		Db:   DB{},
		Jwt:  JWT{},
		Node: Node{},
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
	if db__user := os.Getenv("db__user"); len(db__user) > 0 {
		envConf.Db.User = db__user
	} else if fileConf != nil {
		envConf.Db.User = fileConf.Db.User
	}
	if db__password := os.Getenv("db__password"); len(db__password) > 0 {
		envConf.Db.Password = db__password
	} else if fileConf != nil {
		envConf.Db.Password = fileConf.Db.Password
	}
	if db__host := os.Getenv("db__host"); len(db__host) > 0 {
		envConf.Db.Host = db__host
	} else if fileConf != nil {
		envConf.Db.Host = fileConf.Db.Host
	}
	if db__schema := os.Getenv("db__schema"); len(db__schema) > 0 {
		envConf.Db.Schema = db__schema
	} else if fileConf != nil {
		envConf.Db.Schema = fileConf.Db.Schema
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
	if node__globalname := os.Getenv("node__globalname"); len(node__globalname) > 0 {
		envConf.Node.GlobalName = node__globalname
	} else if fileConf != nil {
		envConf.Node.GlobalName = fileConf.Node.GlobalName
	}
	if node__externaladdr := os.Getenv("node__externaladdr"); len(node__externaladdr) > 0 {
		envConf.Node.ExternalAddr = node__externaladdr
	} else if fileConf != nil {
		envConf.Node.ExternalAddr = fileConf.Node.ExternalAddr
	}
	if node__externalport := os.Getenv("node__externalport"); len(node__externalport) > 0 {
		if port, err := strconv.Atoi(node__externalport); err == nil && port > 0 {
			envConf.Node.ExternalPort = port
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
		if port, err := strconv.Atoi(node__bindport); err == nil && port > 0 {
			envConf.Node.BindPort = port
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

	return envConf
}
