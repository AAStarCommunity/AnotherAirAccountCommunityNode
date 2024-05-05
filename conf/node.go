package conf

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Node struct {
	Genesis      bool
	ExternalAddr string `yaml:"externalAddr"`
	ExternalPort uint16 `yaml:"externalPort"`
	BindAddr     string `yaml:"bindAddr"`
	BindPort     uint16 `yaml:"bindPort"`
	GlobalName   string `yaml:"globalName"`
}

var node *Node

var onceNode sync.Once

func GetNode() *Node {
	onceNode.Do(func() {
		if node == nil {
			j := &getConf().Node
			node = &Node{
				Genesis:      j.Genesis,
				ExternalAddr: j.ExternalAddr,
				ExternalPort: j.ExternalPort,
				BindAddr:     j.BindAddr,
				BindPort:     j.BindPort,
				GlobalName: func() string {
					if j.GlobalName == "" {
						return fmt.Sprintf("aa:%s", uuid.NewString())
					} else {
						return j.GlobalName
					}
				}(),
			}
		}
	})

	return node
}
