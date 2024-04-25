package community

import (
	"another_node/conf"
	"another_node/internal/community/node"
	"another_node/internal/community/storage"
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

// BindAccount binding a web2 account
func BindAccount(hashedAccount string, publicKey *string) error {
	privateKeyValut := "WIP Private Key Vault"

	if publicKey == nil {
		// TODO: auto dispatch a web3 account
		publicKey = new(string)
		*publicKey = "WIP"
		privateKeyValut = "Auto Dispatched Private Key Vault"
	}

	rpcAddress := conf.GetNode().ExternalAddr
	rpcPort := conf.GetNode().ExternalPort
	version := 0

	if err := storage.UpsertMember(hashedAccount, *publicKey, privateKeyValut, rpcAddress, rpcPort, &version); err != nil {
		return err
	} else {
		return community.Node.Broadcast(&node.Payload{
			Account:    hashedAccount,
			PublicKey:  *publicKey,
			RpcAddress: rpcAddress,
			RpcPort:    rpcPort,
			Version:    version,
		})
	}
}
