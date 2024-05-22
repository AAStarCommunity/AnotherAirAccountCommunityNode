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
		payload := storage.Members{
			{
				HashedAccount: hashedAccount,
				PublicKey:     *publicKey,
				RpcAddress:    rpcAddress,
				RpcPort:       rpcPort,
				Version:       version,
			},
		}.Marshal()
		return c.Node.Broadcast(MemberStream, payload)
	}
}

func ListNodes() []string {
	var members []string
	for _, node := range c.Node.Members.Members() {
		members = append(members, node.Name)
	}
	return members
}
