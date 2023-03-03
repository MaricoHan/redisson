package entity

import (
	"time"
)

// Project 项目表
type Project struct {
	Id          uint      `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	Code        string    `gorm:"column:code;type:char(16);comment:项目 Code;NOT NULL" json:"code"`
	ApiKey      string    `gorm:"column:api_key;type:varchar(40);comment:项目 Key;NOT NULL" json:"api_key"`
	ApiSecret   string    `gorm:"column:api_secret;type:char(40);comment:项目密钥;NOT NULL" json:"api_secret"`
	AccessMode  int       `gorm:"column:access_mode;type:tinyint(1);default:1;comment:项目的接入方式 1：托管 2：非托管;NOT NULL" json:"access_mode"`
	ChainId     uint      `gorm:"column:chain_id;type:bigint(20) unsigned;default:0;comment:链 ID;NOT NULL" json:"chain_id"`
	UserId      uint      `gorm:"column:user_id;type:bigint(20) unsigned;default:0;comment:用户 ID;NOT NULL" json:"user_id"`
	Name        string    `gorm:"column:name;type:varchar(255);comment:项目名称;NOT NULL" json:"name"`
	Description string    `gorm:"column:description;type:varchar(255);comment:项目描述;NOT NULL" json:"description"`
	Status      uint      `gorm:"column:status;type:tinyint(4) unsigned;default:1;comment:状态（1.启用 2.禁用 3.注销）;NOT NULL" json:"status"`
	Version     uint      `gorm:"column:version;type:tinyint(4) unsigned;default:0;comment:项目参数版本号,1:V1；2:V2;NOT NULL" json:"version"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
}

const (
	MANAGED = iota + 1
	UNMANAGED
)

var ProjectFields = struct {
	BaseModelFields
	Code        string
	ApiKey      string
	ApiSecret   string
	AccessMode  string
	ChainId     string
	UserId      string
	Name        string
	Description string
	Status      string
	Version     string
}{
	baseModelFields,
	"code",
	"api_key",
	"api_secret",
	"access_mode",
	"chain_id",
	"user_id",
	"name",
	"description",
	"status",
	"version",
}

// 项目参数版本
const (
	Version1 = iota + 1
	Version2
)
