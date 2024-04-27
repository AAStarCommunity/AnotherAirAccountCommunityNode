package node

import (
	"bytes"
	"encoding/gob"
	"log"
)

func (n *Node) listen() {
	for buf := range n.Delegate.DataChannel {
		decoder := gob.NewDecoder(bytes.NewReader(buf))
		payload := Payload{}
		if err := decoder.Decode(&payload); err != nil {
			log.Printf("Failed to decode broadcast data: %v", err)
			continue
		}

		log.Printf("Received broadcast: %v", payload)

		go UpcomingHandler(&payload)
	}
}
