package storage

import (
	"github.com/syndtr/goleveldb/leveldb/util"
)

func GetAllMembers(total uint32) []Member {
	if ins, err := Open(); err != nil {
		return nil
	} else {
		defer ins.Close()
		db := ins.Instance

		members := make([]Member, 0)
		iter := db.NewIterator(&util.Range{
			Start: []byte(MemberPrefix),
		}, nil)
		i := 0
		for iter.Next() {
			if i < int(total) {
				if m, err := Unmarshal(iter.Value()); err != nil {
					return nil
				} else {
					members = append(members, *m)
				}
				i++
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
