package chain

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gorm.io/gorm"
)

type IChainRepo interface {
	QueryChainById(chainId uint64) (chain *entity.Chain, err error)
}

type ChainRepo struct {
	db *gorm.DB
}

func NewChainRepo(db *gorm.DB) *ChainRepo {
	return &ChainRepo{db: db}
}

func (c *ChainRepo) QueryChainById(chainId uint64) (chain entity.Chain, err error) {
	err = c.db.Select("id,code,module").Where("id = ? and status = ?", chainId, 1).Find(&chain).Error
	return chain, err
}
