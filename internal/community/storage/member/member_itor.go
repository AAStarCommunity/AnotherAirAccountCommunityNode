package member_storage

import (
	"another_node/internal/community/storage"

	"github.com/syndtr/goleveldb/leveldb/util"
)

func GetMembers(skip, size uint32) Members {
	if db, err := storage.EnsureOpen(); err != nil {
		return nil
	} else {
		members := make([]Member, 0)
		iter := db.NewIterator(
			util.BytesPrefix([]byte(memberPrefix)),
			nil)
		i := uint32(0)
		for iter.Next() {
			if i >= skip && i < skip+size {
				if m, err := UnmarshalMember(iter.Value()); err != nil {
					return nil
				} else {
					members = append(members, *m)
				}
			}
			i++
		}
		iter.Release()
		err = iter.Error()
		if err != nil {
			return nil
		}

		return members
	}
}
