package conf

import "sync"

type Node struct {
	Genesis bool
}

var node *Node

var onceNode sync.Once

func GetNode() *Node {
	onceNode.Do(func() {
		if node == nil {
			j := &getConf().Node
			node = &Node{
				Genesis: j.Genesis,
			}
		}
	})

	return node
}
