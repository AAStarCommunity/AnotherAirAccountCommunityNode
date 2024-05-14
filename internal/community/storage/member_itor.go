package storage

import (
	"another_node/conf"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func GetAllMembers(skip uint32) []Member {
	if stor, err := conf.GetStorage(); err != nil {
		return nil
	} else {
		if db, err := leveldb.Open(stor, nil); err != nil {
			return nil
		} else {
			defer func() {
				stor.Close()
				db.Close()
			}()

			members := make([]Member, 0)
			iter := db.NewIterator(&util.Range{
				Start: []byte(MemberPrefix),
			}, nil)
			i := 1
			for iter.Next() {
				if i >= int(skip) {
					if m, err := Unmarshal(iter.Value()); err != nil {
						return nil
					} else {
						members = append(members, *m)
					}
				}
			}
			iter.Release()
			err = iter.Error()
			if err != nil {
				return nil
			}

			return members
		}
	}
}
