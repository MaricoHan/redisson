package entity

import (
	. "gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
)

// 权限表
type Permission struct {
	Id       uint   `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:权限id" json:"id"`
	Path     string `gorm:"column:path;type:varchar(50);comment:请求路径;NOT NULL" json:"path"`
	Method   string `gorm:"column:method;type:varchar(50);comment:请求方法;NOT NULL" json:"method"`
	Action   int    `gorm:"column:action;type:tinyint(1);comment:操作 1:允许 2:拒绝;NOT NULL" json:"action"`
	Priority int    `gorm:"column:priority;type:int(11);comment:优先级(值越小优先级越高);NOT NULL" json:"priority"`
}

var PermissoinFields = struct {
	ID        string
	Path      string
	Method    string
	Action    string
	Priority  string
	CreatedAt string
	UpdatedAt string
}{
	ID:        MysqlPermissoinTable + ".id",
	Path:      MysqlPermissoinTable + ".path",
	Method:    MysqlPermissoinTable + ".method",
	Action:    MysqlPermissoinTable + ".action",
	Priority:  MysqlPermissoinTable + ".priority",
	CreatedAt: MysqlPermissoinTable + ".created_at",
	UpdatedAt: MysqlPermissoinTable + ".updated_at",
}
