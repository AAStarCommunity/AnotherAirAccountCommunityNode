package node

import (
	"another_node/conf"

	"github.com/hashicorp/memberlist"
)

func init() {

	list, err := memberlist.Create(memberlist.DefaultLocalConfig())

	node = &Node{
		Members: list,
	}

	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	if conf.GetNode().Genesis {
		if err := node.Genesis(); err != nil {
			panic("Failed to start genesis node: " + err.Error())
		}
	} else {
		if err := node.Join(); err != nil {
			panic("Failed to join existing cluster: " + err.Error())
		}
	}
}
