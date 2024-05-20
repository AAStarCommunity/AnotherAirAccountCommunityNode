package node

import (
	"another_node/conf"
	"log"
)

// Broadcast sends data to all nodes in the cluster
func (n *Node) Broadcast(data []byte) error {
	me := conf.GetNode().GlobalName
	if len(n.Delegate.Broadcasts) < n.Delegate.BroadcastCap {
		n.Delegate.Broadcasts = append(n.Delegate.Broadcasts, data)
		log.Printf(me+": Broadcasted: %v", data)
	} else {
		log.Printf(me+": Broadcasted data is full, dropped: %v", data)
	}

	return nil
}
