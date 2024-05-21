package node

import (
	"log"
)

// Broadcast sends data to all nodes in the cluster
func (n *Node) Broadcast(data []byte) error {
	if addr, err := getAddr(); err != nil {
		return err
	} else {
		if len(n.Delegate.Broadcasts) < n.Delegate.BroadcastCap {
			n.Delegate.Broadcasts = append(n.Delegate.Broadcasts, data)
			log.Printf("%s: Broadcasted: %v", addr, data)
		} else {
			log.Printf("%s: Broadcasted data is full, dropped: %v", addr, data)
		}

		return nil
	}
}
