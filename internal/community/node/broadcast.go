package node

import (
	"log"
)

// Broadcast sends data to all nodes in the cluster
func (n *Node) Broadcast(protocol uint8, data []byte) error {
	if addr, err := getAddr(); err != nil {
		return err
	} else {
		if len(n.Delegate.Broadcasts) < n.Delegate.BroadcastCap {
			m := []byte{protocol}
			m = append(m, data...)
			n.Delegate.Broadcasts = append(n.Delegate.Broadcasts, m)
			log.Printf("%s: Protocol: %d Broadcasted: %v", addr, protocol, data)
		} else {
			log.Printf("%s: Broadcasted data is full, dropped: %v", addr, data)
		}

		return nil
	}
}
