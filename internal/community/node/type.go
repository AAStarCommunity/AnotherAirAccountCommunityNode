package node

import "github.com/hashicorp/memberlist"

type Node struct {
	Members *memberlist.Memberlist
}

var node *Node

func (n *Node) Genesis() error {
	// start the genesis node
	return nil
}

func (n *Node) Join() error {
	// join the existing cluster
	return nil
}
