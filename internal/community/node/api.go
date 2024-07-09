package node

import (
	"another_node/conf"
	member_storage "another_node/internal/community/storage/member"
	"another_node/internal/web_server/pkg/request"
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

	if err := member_storage.UpsertMember(hashedAccount, *publicKey, privateKeyValut, rpcAddress, rpcPort, &version); err != nil {
		return err
	} else {
		payload := member_storage.Members{
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
func Sign(req request.Sign) error {
	return nil
}
func ListNodes() []string {
	var members []string
	for _, node := range c.Node.Members.Members() {
		members = append(members, node.Name)
	}
	return members
}
