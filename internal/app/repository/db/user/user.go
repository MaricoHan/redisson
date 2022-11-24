package user

import (
	"gorm.io/gorm"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
)

type IUserRepo interface {
	GetUser(id uint64) (user entity.User, err error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (p *UserRepo) GetUser(id uint64) (user entity.User, err error) {
	err = p.db.Select("id,code").Where("id=?", id).Find(&user).Error
	return user, err
}
