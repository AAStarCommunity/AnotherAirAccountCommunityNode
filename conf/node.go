package conf

import (
	"sync"
)

type Node struct {
	Standalone   bool
	Genesis      bool
	ExternalAddr string `yaml:"externalAddr"`
	ExternalPort uint16 `yaml:"externalPort"`
	BindAddr     string `yaml:"bindAddr"`
	BindPort     uint16 `yaml:"bindPort"`
}

var node *Node

var onceNode sync.Once

func GetNode() *Node {
	onceNode.Do(func() {
		if node == nil {
			j := &getConf().Node
			node = &Node{
				Standalone:   j.Standalone,
				Genesis:      j.Genesis,
				ExternalAddr: j.ExternalAddr,
				ExternalPort: j.ExternalPort,
				BindAddr:     j.BindAddr,
				BindPort:     j.BindPort,
			}
		}
	})

	return node
}
