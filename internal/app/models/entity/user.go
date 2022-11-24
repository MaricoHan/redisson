package entity

import "time"

// User 用户表
type User struct {
	Id               uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	Code             string    `gorm:"column:code;type:char(16);comment:用户 Code;NOT NULL" json:"code"`
	Amount           float64   `gorm:"column:amount;type:decimal(65,4);default:0.0000;comment:人民币余金额" json:"amount"`
	InvoicableAmount float64   `gorm:"column:invoicable_amount;type:decimal(20,4);default:0.0000;comment:可开票金额;NOT NULL" json:"invoicable_amount"`
	UserName         string    `gorm:"column:user_name;type:varchar(20);comment:用户名" json:"user_name"`
	Password         string    `gorm:"column:password;type:char(40);comment:登录密码（SHA1( MD5 (明文密码) + 盐值）;NOT NULL" json:"password"`
	Salt             string    `gorm:"column:salt;type:varchar(10);comment:盐值;NOT NULL" json:"salt"`
	Icon             string    `gorm:"column:icon;type:varchar(255);comment:图像链接;NOT NULL" json:"icon"`
	Type             uint      `gorm:"column:type;type:tinyint(4) unsigned;default:1;comment:用户类型（1：个人用户（默认）2：企业用户）;NOT NULL" json:"type"`
	PhoneNumber      string    `gorm:"column:phone_number;type:varchar(20);comment:手机号" json:"phone_number"`
	Email            string    `gorm:"column:email;type:varchar(64);comment:电子邮箱" json:"email"`
	Description      string    `gorm:"column:description;type:text;comment:个人介绍" json:"description"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
}
