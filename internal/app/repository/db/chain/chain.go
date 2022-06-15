package chain

import (
	"errors"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gorm.io/gorm"
)

type IChainRepo interface {
	CreateChainTable() error
	InsertChain(chain entity.Chain) error
	QueryChainById(chainId uint64) (chain *entity.Chain, err error)
}

type ChainRepo struct {
	db *gorm.DB
}

func NewChainRepo(db *gorm.DB) *ChainRepo {
	return &ChainRepo{db: db}
}

func (c *ChainRepo) CreateChainTable() error {
	if !c.db.Migrator().HasTable(&entity.Chain{}) {
		err := c.db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(&entity.Chain{})
		if err != nil {
			return err
		}
		if c.db.Migrator().HasTable(&entity.Chain{}) {
			return nil
		} else {
			return errors.New("create chain table failed")
		}
	} else {
		return nil
	}
}

func (c *ChainRepo) InsertChain(chain entity.Chain) error {
	return c.db.Create(&chain).Error
}

func (c *ChainRepo) QueryChainById(chainId uint64) (chain entity.Chain, err error) {
	err = c.db.Where("id = ? and status = ?", chainId, 1).Find(&chain).Error
	return chain, err

}
