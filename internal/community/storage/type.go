package storage

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeletedBy string         `gorm:"size:64, null"`
}

// gorm的数据表结构定义
type Member struct {
	Model
	HashedAccount   string `gorm:"uniqueIndex"`
	PublicKey       string `gorm:"size:1024"`
	PrivateKeyVault string `gorm:"size:1024, null"`
}
