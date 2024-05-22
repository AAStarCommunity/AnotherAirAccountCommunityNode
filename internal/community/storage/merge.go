package storage

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
		if db, err := EnsureOpen(); err == nil {
			db.Put([]byte(NodeKey(&node)), node.Marshal(), nil)
		}
	}
}
