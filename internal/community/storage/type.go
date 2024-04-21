package storage

import (
	"os"
	"time"

	"gorm.io/gorm"
)

var rpcAddress string

type Model struct {
	ID        uint           `gorm:"column:id; primarykey"`
	CreatedAt time.Time      `gorm:"column:created_at; autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at; autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at; index"`
	DeletedBy string         `gorm:"column:deleted_by; type: varchar(1024); null"`
	Version   uint           `gorm:"column:version; default:0;"`
}

func init() {
	rpcAddress = os.Getenv("RPC_ADDRESS")
}
