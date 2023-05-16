package entity

import (
	. "gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"time"
)

// 服务权限表
type ServiceXPermissoin struct {
	Id           uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	ServiceId    int       `gorm:"column:service_id;type:int(11);comment:服务id;NOT NULL" json:"service_id"`
	PermissionId int       `gorm:"column:permission_id;type:int(11);comment:权限id;NOT NULL" json:"permission_id"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}

var ServiceXPermissoinFields = struct {
	ID           string
	ServiceId    string
	PermissionId string
}{
	ID:           MysqlServiceXPermissoinTable + ".id",
	ServiceId:    MysqlServiceXPermissoinTable + ".service_id",
	PermissionId: MysqlServiceXPermissoinTable + ".permission_id",
}
