package storage

import (
	"another_node/conf"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type Db struct {
	stor     storage.Storage
	Instance *leveldb.DB
	lc       sync.Mutex
}

func (d *Db) Close() {
	d.stor.Close()
	d.Instance.Close()
	d.lc.Unlock()
}

func Open() (*Db, error) {
	if stor, err := conf.GetStorage(); err != nil {
		return nil, err
	} else {
		if db, err := leveldb.Open(stor, nil); err != nil {
			return nil, err
		} else {
			ins := &Db{
				stor:     stor,
				Instance: db,
			}

			ins.lc.Lock()
			return ins, nil
		}
	}
}
