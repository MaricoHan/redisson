package entity

import (
	"time"
)

// 权限表
type Permission struct {
	Id        uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:权限id" json:"id"`
	Path      string    `gorm:"column:path;type:varchar(50);comment:请求路径;NOT NULL" json:"path"`
	Method    string    `gorm:"column:method;type:varchar(50);comment:请求方法;NOT NULL" json:"method"`
	Action    int       `gorm:"column:action;type:tinyint(1);comment:操作 1:允许 2:拒绝;NOT NULL" json:"action"`
	Priority  int       `gorm:"column:priority;type:int(11);comment:优先级(值越小优先级越高);NOT NULL" json:"priority"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}
