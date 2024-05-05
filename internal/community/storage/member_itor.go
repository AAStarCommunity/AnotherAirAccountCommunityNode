package storage

import (
	"another_node/conf"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func GetAllMembers() ([]Member, error) {
	if stor, err := conf.GetStorage(); err != nil {
		return nil, err
	} else {
		if db, err := leveldb.Open(stor, nil); err != nil {
			return nil, err
		} else {
			defer func() {
				stor.Close()
				db.Close()
			}()

			members := make([]Member, 0)
			iter := db.NewIterator(&util.Range{
				Start: []byte(MemberPrefix),
			}, nil)
			for iter.Next() {
				if m, err := Unmarshal(iter.Value()); err != nil {
					return nil, err
				} else {
					members = append(members, *m)
				}
			}
			iter.Release()
			err = iter.Error()
			if err != nil {
				return nil, err
			}

			return members, nil
		}
	}
}
