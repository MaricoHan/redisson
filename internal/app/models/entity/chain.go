package entity

import (
	"github.com/volatiletech/sqlboiler/types"
	constant "gitlab.bianjie.ai/avata/open-api/internal/app/models"
	"time"
)

type Chain struct {
	ID        uint64            `gorm:"column:id;type:bigint(20) auto_increment;primaryKey;not null"`
	Code      string            `gorm:"column:code;type:varchar(15);not null;unique;default:'';comment:链code" json:"code"`
	Name      string            `gorm:"column:name;type:varchar(255);not null;default:'';comment:链名称" json:"name"`
	Module    string            `gorm:"column:module;type:varchar(255);not null;default:'';comment:链模块" json:"module"`
	GasPrice  types.NullDecimal `gorm:"column:amount;type:decimal(10,8);default:0;comment:换算对比" json:"amount"`
	Status    int               `gorm:"column:status;type:tinyint(3);not null;default:1;comment:链状态,1为正常"`
	CreatedAt time.Time         `gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time         `gorm:"column:updated_at;type:datetime"`
}

func (*Chain) TableName() string {
	return constant.TableChain
}
