package node

import (
	"another_node/conf"
	"log"

	"github.com/hashicorp/memberlist"
)

func New() (*Node, error) {
	confNode := conf.GetNode()
	conf := memberlist.DefaultWANConfig()
	conf.Name = confNode.GlobalName
	conf.AdvertiseAddr = confNode.ExternalAddr
	conf.AdvertisePort = confNode.ExternalPort
	conf.BindAddr = confNode.BindAddr
	conf.BindPort = confNode.BindPort

	var err error
	var list *memberlist.Memberlist
	dataChan := make(chan []byte, 10)

	if list, err = memberlist.Create(conf); err == nil {

		if !confNode.Genesis {
			exists := []string{
				"192.168.1.6:7947", // TODO: replace with the genesis node address on chain
			}

			if _, err := list.Join(exists); err != nil {
				log.Fatalf("Failed to join cluster: %v", err)
				return nil, err
			}
		}

		delegate := &CommunityDelegate{
			DataChannel:  dataChan,
			Broadcasts:   make([][]byte, 0),
			BroadcastCap: 10,
		}

		node := &Node{
			Members:  list,
			Delegate: delegate,
		}

		go node.listen()

		return node, nil
	}

	return nil, err
}
