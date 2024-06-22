package member_storage

import "another_node/internal/community/storage"

func MergeRemoteAccounts(accounts []Member) error {
	for _, account := range accounts {
		if err := UpsertMember(
			account.HashedAccount,
			account.PublicKey,
			account.PrivateKeyVault,
			account.RpcAddress,
			account.RpcPort,
			&account.Version); err != nil {
			return err
		}
	}
	return nil
}

func MergeRemoteAddr(buf []byte) {
	nodes := unmarshalNodes(buf)
	for _, node := range nodes {
		if db, err := storage.EnsureOpen(); err == nil {
			db.Put([]byte(NodeKey(&node)), node.Marshal(), nil)
		}
	}
}
