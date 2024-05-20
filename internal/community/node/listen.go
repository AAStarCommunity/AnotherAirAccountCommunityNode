package node

import (
	"another_node/conf"
	"log"
)

func (n *Node) listen() {
	me := conf.GetNode().GlobalName
	for buf := range n.Delegate.DataChannel {
		log.Printf(me+": Received broadcast: %v", buf)
		go UpcomingHandler(buf)
	}
}
