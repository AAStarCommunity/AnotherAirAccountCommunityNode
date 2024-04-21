package storage

import "another_node/conf"

func prepare_test() {
	db := conf.GetDbClient()

	db.AutoMigrate(&Member{})
}
