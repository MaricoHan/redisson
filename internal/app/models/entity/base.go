package entity

import "time"

type BaseModel struct {
	ID        uint64    `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"<-:false"`
	UpdatedAt time.Time `gorm:"<-:false"`
}

type BaseModelFields = struct {
	ID        string
	CreatedAt string
	UpdatedAt string
}

var baseModelFields = BaseModelFields{
	"id",
	"created_at",
	"updated_at",
}
