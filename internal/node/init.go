package node

import (
	"fmt"

	"github.com/hashicorp/memberlist"
)

var members *memberlist.Memberlist

func init() {
	list, err := memberlist.Create(memberlist.DefaultLocalConfig())
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}
	// Join an existing cluster by specifying at least one known member.
	// TODO: This should be read from a smart contract from Layer2 chain
	_, err = list.Join([]string{"1.2.3.4"})
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}

	// Ask for members of the cluster
	for _, member := range list.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}

	node = &Node{
		Members: list,
	}
}
