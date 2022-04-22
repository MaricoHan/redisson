package entity

import (
	constant "gitlab.bianjie.ai/avata/open-api/internal/app/models"
	"time"
)

type Project struct {
	ID          uint64    `gorm:"column:id;type:bigint(20) auto_increment;primaryKey;not null"`
	ApiKey      string    `gorm:"column:api_key;type:varchar(255);not null;default:'';comment:api_key" json:"api_key"`
	ApiSecret   string    `gorm:"column:api_secret;type:varchar(255);not null;default:'';comment:api_secret" json:"api_secret"`
	Name        string    `gorm:"column:name;type:varchar(255);binary;not null;default:'';comment:项目名称" json:"name"`
	Description string    `gorm:"column:description;type:varchar(255);not null;default:'';comment:项目描述" json:"description"`
	ChainID     uint64    `gorm:"column:chain_id;type:bigint(20);not null;default:0;comment:链id" json:"chain_id"`
	UserID      uint64    `gorm:"column:user_id;type:bigint(20);not null;default:0;" json:"user_id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime"`
}

func (*Project) TableName() string {
	return constant.TableProject
}
