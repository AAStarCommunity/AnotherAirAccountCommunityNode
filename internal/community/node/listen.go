package node

import (
	"bytes"
	"encoding/gob"
	"log"
)

func (n *Node) listen() {
	for {
		select {
		case buf := <-n.Delegate.DataChannel:
			decoder := gob.NewDecoder(bytes.NewReader(buf))
			payload := Payload{}
			if err := decoder.Decode(&payload); err != nil {
				log.Printf("Failed to decode broadcast data: %v", err)
				continue
			}

			go UpcomingHandler(&payload)
		}
	}
}
