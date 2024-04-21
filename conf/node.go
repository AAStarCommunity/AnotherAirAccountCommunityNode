package conf

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Node struct {
	Genesis         bool
	ExternalAddress string
	GlobalName      string
}

var node *Node

var onceNode sync.Once

func GetNode() *Node {
	onceNode.Do(func() {
		if node == nil {
			j := &getConf().Node
			node = &Node{
				Genesis:         j.Genesis,
				ExternalAddress: j.ExternalAddress,
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
