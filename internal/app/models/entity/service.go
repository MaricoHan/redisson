package entity

import (
	"time"
)

const (
	ServiceTypeWallet = 1 // 钱包服务
	ServiceTypeNS     = 2 // 域名服务
)

// Services 服务类型表
type Services struct {
	Id        uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	Name      string    `gorm:"column:name;type:char(10);comment:服务名称;NOT NULL" json:"name"`
	Type      uint      `gorm:"column:type;type:tinyint(4) unsigned;default:0;comment:服务类型, 1: 钱包;NOT NULL" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}
