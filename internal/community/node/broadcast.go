package node

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Broadcast sends data to all nodes in the cluster
func (n *Node) Broadcast(data *Payload) error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		log.Printf("Failed to serialize data: %v", err)
		return err
	}

	if len(n.Delegate.Broadcasts) < n.Delegate.BroadcastCap {
		n.Delegate.Broadcasts = append(n.Delegate.Broadcasts, buf.Bytes())
		log.Printf("Broadcasted: %v", data)
	} else {
		log.Printf("Broadcasted data is full, dropped: %v", data)
	}

	return nil
}
