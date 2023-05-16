package entity

import (
	. "gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"time"
)

// Services 服务类型表
type Services struct {
	Id       uint   `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	ParentId int    `gorm:"column:parent_id;type:int(11);default:0;comment:父服务id;NOT NULL" json:"parent_id"`
	Name     string `gorm:"column:name;type:char(10);comment:服务名称;NOT NULL" json:"name"`
	//Type      uint      `gorm:"column:type;type:tinyint(4) unsigned;default:0;comment:服务类型, 1: 钱包;NOT NULL" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}

var ServiceFields = struct {
	ID       string
	ParentId string
	Name     string
}{
	ID:       MysqlServicesTable + ".id",
	ParentId: MysqlServicesTable + ".parent_id",
	Name:     MysqlServicesTable + ".name",
}
