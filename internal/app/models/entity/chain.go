package entity

import "time"

// 链表
type Chain struct {
	Id                uint      `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	Code              string    `gorm:"column:code;type:varchar(15);comment:链 Code,如: wenchangchain;NOT NULL" json:"code"`
	Name              string    `gorm:"column:name;type:varchar(255);comment:名称;NOT NULL" json:"name"`
	Module            string    `gorm:"column:module;type:varchar(100);comment:模块：native, ddc;NOT NULL" json:"module"`
	GasPrice          float64   `gorm:"column:gas_price;type:decimal(10,8);default:0.00000000;comment:换算对比;NOT NULL" json:"gas_price"`
	Status            int       `gorm:"column:status;type:tinyint(4) unsigned;default:0;comment:链状态（1: 有效 2：无效）;NOT NULL" json:"status"`
	Description       string    `gorm:"column:description;type:varchar(500);not null;default:''"`
	ModuleDescription string    `gorm:"column:module_description;type:varchar(500);not null;default:''"`
	UpdatedAt         time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
	CreatedAt         time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
}
