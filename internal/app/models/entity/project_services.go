package entity

import "time"

type ProjectXServices struct {
	Id        uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	ProjectId uint      `gorm:"column:project_id;type:int(11) unsigned;default:0;comment:项目 ID;NOT NULL" json:"project_id"`
	ServiceId uint      `gorm:"column:service_id;type:int(11) unsigned;default:0;comment:服务 ID;NOT NULL" json:"service_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}
