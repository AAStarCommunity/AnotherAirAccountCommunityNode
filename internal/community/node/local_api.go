package node

import (
	"another_node/conf"
	"another_node/internal/community/storage"
)

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
	version := uint32(0)

	if err := storage.UpsertMember(hashedAccount, *publicKey, privateKeyValut, rpcAddress, rpcPort, &version); err != nil {
		return err
	} else {
		return c.Node.Broadcast(&Payload{
			Account:    hashedAccount,
			PublicKey:  *publicKey,
			RpcAddress: rpcAddress,
			RpcPort:    rpcPort,
			Version:    version,
		})
	}
}

func ListNodes() []string {
	var members []string
	for _, node := range c.Node.Members.Members() {
		members = append(members, node.Name)
	}
	return members
}

func Broadcast(payload *Payload) error {
	return c.Node.Broadcast(payload)
}
