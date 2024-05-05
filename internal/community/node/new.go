package node

import (
	"another_node/conf"
	"log"
	"strings"

	"github.com/hashicorp/memberlist"
)

func New(listen *uint16, globalName *string, entrypoints *string, genesis *bool) (*Node, error) {

	confNode := conf.GetNode()
	if genesis != nil {
		confNode.Genesis = *genesis
	}
	if listen != nil && *listen > 0 {
		confNode.ExternalPort = *listen
	}
	if globalName != nil && len(*globalName) > 0 {
		confNode.GlobalName = *globalName
	}

	delegate := &CommunityDelegate{
		DataChannel:  make(chan []byte, 10),
		Broadcasts:   make([][]byte, 0),
		BroadcastCap: 10,
	}

	conf := memberlist.DefaultWANConfig()
	conf.Name = confNode.GlobalName
	conf.AdvertiseAddr = confNode.ExternalAddr
	conf.AdvertisePort = int(confNode.ExternalPort)
	conf.BindAddr = confNode.BindAddr
	conf.BindPort = int(confNode.BindPort)
	conf.Delegate = delegate

	entrypointNodeAddr := []string{
		"192.168.1.6:7947", // TODO: replace with the genesis node address on chain
	}
	if entrypoints != nil && len(*entrypoints) > 0 {
		entrypointNodeAddr = strings.Split(*entrypoints, ",")
	}

	var err error
	var list *memberlist.Memberlist

	if list, err = memberlist.Create(conf); err == nil {

		if !confNode.Genesis {
			exists := entrypointNodeAddr

			if _, err := list.Join(exists); err != nil {
				log.Fatalf("Failed to join cluster: %v", err)
				return nil, err
			}
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
