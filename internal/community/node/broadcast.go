package node

import (
	"bytes"
	"encoding/gob"
	"log"
)

// broadcast sends data to all nodes in the cluster
func (n *Node) broadcast(data *Payload) error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		log.Printf("Failed to serialize data: %v", err)
		return err
	}

	if len(n.Delegate.Broadcasts) < n.Delegate.BroadcastCap {
		n.Delegate.Broadcasts = append(n.Delegate.Broadcasts, buf.Bytes())
	}

	return nil
}
