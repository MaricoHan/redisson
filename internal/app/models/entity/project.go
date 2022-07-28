package entity

import (
	"time"
)

// 项目表
type Project struct {
	Id          uint      `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	Code        string    `gorm:"column:code;type:char(16);comment:项目 Code;NOT NULL" json:"code"`
	ApiKey      string    `gorm:"column:api_key;type:varchar(40);comment:项目 Key;NOT NULL" json:"api_key"`
	ApiSecret   string    `gorm:"column:api_secret;type:char(40);comment:项目密钥;NOT NULL" json:"api_secret"`
	ChainId     uint      `gorm:"column:chain_id;type:bigint(20) unsigned;default:0;comment:链 ID;NOT NULL" json:"chain_id"`
	UserId      uint      `gorm:"column:user_id;type:bigint(20) unsigned;default:0;comment:用户 ID;NOT NULL" json:"user_id"`
	Name        string    `gorm:"column:name;type:varchar(255);comment:项目名称;NOT NULL" json:"name"`
	Description string    `gorm:"column:description;type:varchar(255);comment:项目描述;NOT NULL" json:"description"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
}
