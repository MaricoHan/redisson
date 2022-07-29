package entity

// 链表
type Chain struct {
	Id     uint   `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	Code   string `gorm:"column:code;type:varchar(15);comment:链 Code,如: wenchangchain;NOT NULL" json:"code"`
	Name   string `gorm:"column:name;type:varchar(255);comment:名称;NOT NULL" json:"name"`
	Module string `gorm:"column:module;type:varchar(100);comment:模块：native, ddc;NOT NULL" json:"module"`
	Status int    `gorm:"column:status;type:tinyint(4) unsigned;default:0;comment:链状态（1: 有效 2：无效）;NOT NULL" json:"status"`
}
