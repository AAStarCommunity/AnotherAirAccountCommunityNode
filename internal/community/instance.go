package community

import (
	"another_node/conf"
	"another_node/internal/community/node"
	"another_node/internal/community/storage"
	"fmt"
)

type Community struct {
	Node *node.Node
}

var community *Community

func New(n *node.Node) {
	community = &Community{
		Node: n,
	}
}

func BindAccount(hashedAccount string, publicKey *string) error {
	if publicKey == nil {
		// TODO: auto dispatch a web3 account
		publicKey = new(string)
		*publicKey = "WIP"
	}

	rpcAddress := new(string)
	*rpcAddress = fmt.Sprintf("%s:%d", conf.GetNode().ExternalAddr, conf.GetNode().ExternalPort)

	if err := storage.UpsertMember(hashedAccount, publicKey, nil, rpcAddress, 0); err != nil {
		return err
	} else {
		return community.Node.Broadcast(&node.Payload{
			Account:    hashedAccount,
			PublicKey:  *publicKey,
			RpcAddress: *rpcAddress,
			Version:    0,
		})
	}
}
