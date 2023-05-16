package user

import (
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
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
	err = p.db.Select(entity.UserFields.ID, entity.UserFields.Code).
		Where(entity.UserFields.ID, id).
		Where(entity.UserFields.IsDeleted, constant.IsNotDelete).
		Find(&user).Error

	return
}
