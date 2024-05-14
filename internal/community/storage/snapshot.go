package storage

import (
	"another_node/conf"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

const Version uint8 = 1

type Snapshot struct {
	Version      uint8
	TotalMembers uint32
	HashedDigest []byte
}

func (s *Snapshot) Digest() *Snapshot {
	s.HashedDigest = []byte{1, 2, 3}
	return s
}

func (s *Snapshot) Marshal() []byte {
	sizeOfSnapshot := 1 + 4 + len(s.HashedDigest)
	buf := make([]byte, sizeOfSnapshot)

	offset := 0
	buf[offset] = s.Version
	offset += 1
	copy(buf[offset:offset+4], uintToBytes(s.TotalMembers))
	offset += 4
	copy(buf[offset:], s.HashedDigest)
	return buf
}

func UnmarshalSnapshot(data []byte) (*Snapshot, error) {
	if len(data) < 1+4 {
		return nil, errors.New("invalid snapshot data")
	}

	s := &Snapshot{}
	offset := 0
	s.Version = data[offset]
	offset += 1
	s.TotalMembers = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4
	s.HashedDigest = data[offset:]
	return s, nil
}

const SnapshotKey = "snapshot"

// GetSnapshot returns the current node state
// node state represent the members of the community snapshot
func GetSnapshot() ([]byte, error) {
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

			if ss, err := db.Get([]byte(SnapshotKey), nil); err != nil {
				if errors.Is(err, leveldb.ErrNotFound) {
					return nil, nil
				}
				return nil, err
			} else {
				return ss, nil
			}
		}
	}
}

func memberCounter() (total uint32) {
	if stor, err := conf.GetStorage(); err != nil {
		return
	} else {
		if db, err := leveldb.Open(stor, nil); err != nil {
			return
		} else {
			defer func() {
				stor.Close()
				db.Close()
			}()

			iter := db.NewIterator(nil, nil)
			total = 0
			for iter.Next() {
				total++
			}
			iter.Release()
			err = iter.Error()
			if err != nil {
				return
			}

			return
		}
	}
}

func ScheduleSnapshot() {
	t := time.NewTicker(10 * time.Second)
	for range t.C {
		if err := updateSnapshot(); err != nil {
			// log error
			fmt.Println("error updating snapshot")
		}
	}
}

func updateSnapshot() error {
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

			ss := &Snapshot{
				Version:      Version,
				TotalMembers: memberCounter(),
			}
			marshal := ss.Digest().Marshal()
			return db.Put([]byte(SnapshotKey), marshal, nil)
		}
	}
}
