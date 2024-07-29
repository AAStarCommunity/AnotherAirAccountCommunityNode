package storage

import (
	"another_node/conf"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var db *Db
var er error

type Db struct {
	stor    storage.Storage
	LevelDB *leveldb.DB
}

func Close() {
	db.stor.Close()
	db.LevelDB.Close()
}

var mutex sync.Mutex

func EnsureOpen() (*leveldb.DB, error) {
	if db != nil {
		return db.LevelDB, nil
	} else {
		mutex.Lock()
		defer mutex.Unlock()
		if db == nil {
			if stor, err := conf.GetStorage(); err == nil {
				if dx, err := leveldb.Open(stor, nil); err == nil {
					db = &Db{
						stor:    stor,
						LevelDB: dx,
					}
				} else {
					er = err
				}
			} else {
				er = err
			}
		}
	}
	return db.LevelDB, er
}
