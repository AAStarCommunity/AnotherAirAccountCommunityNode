package storage

import (
	"another_node/conf"

	"github.com/syndtr/goleveldb/leveldb"
)

func newMemberIndex(member *Member) error {
	if stor, err := conf.GetStorage(); err != nil {
		return err
	} else {
		if db, err := leveldb.Open(stor, nil); err != nil {
			return err
		} else {
			defer func() {
				stor.Close()
				db.Close()
			}()

			if err := db.Put([]byte(memberIndexKey()), []byte(member.HashedAccount), nil); err != nil {
				return err
			}
		}
	}

	return nil
}
