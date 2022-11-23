package entity

import "time"

// ServiceRedirectUrl 服务重定向网关地址表
type ServiceRedirectUrl struct {
	Id        uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	ProjectId uint64    `gorm:"column:project_id;type:bigint(20) unsigned;default:0;comment:项目 ID;NOT NULL" json:"project_id"`
	Url       string    `gorm:"column:url;type:varchar(255);comment:网关地址;NOT NULL" json:"url"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}
