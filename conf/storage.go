package conf

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb/storage"
)

func GetStorage() (storage.Storage, error) {
	if os.Getenv("UnitTest") == "1" {
		return storage.NewMemStorage(), nil
	} else {
		return storage.OpenFile(getConf().Storage, false)
	}
}
