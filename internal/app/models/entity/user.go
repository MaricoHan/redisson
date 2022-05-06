package entity

import (
	"github.com/volatiletech/sqlboiler/v4/types"
	constant "gitlab.bianjie.ai/avata/open-api/internal/app/models"
	"gorm.io/gorm"
)

type (
	LoginType int //登录类型
)

const (
	UserTypeAccount LoginType = 1 //账号密码
	UserTypePhone   LoginType = 2 //手机号，验证码
)

// User 用户信息表
type User struct {
	gorm.Model
	Amount          types.NullDecimal `gorm:"column:amount;type:decimal(20,2);default:0" json:"amount"`
	UserName        string            `gorm:"column:user_name;type:varchar(20);not null;default:'';uniqueIndex:uk_id_username_phone_number_email,priority:2" json:"name"`
	Password        string            `gorm:"column:password;type:varchar(50);not null;default:''" json:"password,omitempty"`
	Salt            string            `gorm:"column:salt;type:varchar(10);not null;default:''" json:"salt,omitempty"`
	Icon            string            `gorm:"column:icon;type:varchar(255);default:''" json:"icon"`
	Type            int               `gorm:"column:type;type:tinyint(1);not null;default:1" json:"type"`
	Phonconstantber string            `gorm:"column:phone_number;type:char(11);not null;default:'';uniqueIndex:uk_id_username_phone_number_email,priority:3" json:"phone"`
	Email           string            `gorm:"column:email;type:varchar(64);not null;default:'';uniqueIndex:uk_id_username_phone_number_email,priority:4" json:"email"`
	Description     string            `gorm:"column:description;type:text" json:"description"`
}

func (*User) TableName() string {
	return constant.TableUser
}
