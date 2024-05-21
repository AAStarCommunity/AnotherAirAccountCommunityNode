package node

import (
	"log"
)

func (n *Node) listen() {
	if addr, err := getAddr(); err == nil {
		for buf := range n.Delegate.DataChannel {
			log.Printf("%s: Received broadcast: %v", addr, buf)
			go UpcomingHandler(buf)
		}
	}
}
