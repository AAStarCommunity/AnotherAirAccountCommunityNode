package storage

import (
	"github.com/syndtr/goleveldb/leveldb/util"
)

func GetAllMembers(skip uint32) []Member {
	if ins, err := Open(); err != nil {
		return nil
	} else {
		defer ins.Close()
		db := ins.Instance

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
