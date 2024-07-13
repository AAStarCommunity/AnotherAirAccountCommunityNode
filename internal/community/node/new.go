package node

import (
	"another_node/conf"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/memberlist"
)

type Community struct {
	Node *Node
}

var c *Community
var entrypointNodeAddr []string

func New(listen *uint16, globalName *string, entrypoints *string, genesis *bool) (*Community, error) {

	confNode := conf.GetNode()
	if confNode.Standalone {
		fmt.Println("Node is running in standalone mode")
		return &Community{
			Node: &Node{},
		}, nil
	} else {
		if genesis != nil {
			confNode.Genesis = *genesis
		}
		if listen != nil && *listen > 0 {
			confNode.ExternalPort = *listen
		}

		delegate := &CommunityDelegate{
			DataChannel:  make(chan []byte, 10),
			Broadcasts:   make([][]byte, 0),
			BroadcastCap: 10,
		}

		conf := memberlist.DefaultWANConfig()
		conf.Name = func() string {
			if addr, err := getAddr(); err != nil {
				panic(err)
			} else {
				return string(addr)
			}
		}()
		conf.AdvertiseAddr = confNode.ExternalAddr
		conf.AdvertisePort = int(confNode.ExternalPort)
		conf.BindAddr = confNode.BindAddr
		conf.BindPort = int(confNode.BindPort)
		conf.Delegate = delegate
		conf.PushPullInterval = time.Second

		entrypointNodeAddr = []string{
			"192.168.1.6:7947", // TODO: replace with the genesis node address on chain
		}
		if entrypoints != nil && len(*entrypoints) > 0 {
			entrypointNodeAddr = strings.Split(*entrypoints, ",")
		}

		var err error
		var list *memberlist.Memberlist
		fmt.Printf("ready to create memberlist with conf: %v\n", conf)
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

			c = &Community{
				Node: node,
			}

			return c, nil
		}

		return nil, err
	}
}
