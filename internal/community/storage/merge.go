package storage

func MergeRemoteAccounts(accounts []Member) error {
	for _, account := range accounts {
		if err := MergeRemoteMember(&account); err != nil {
			return err
		}
	}
	return nil
}
