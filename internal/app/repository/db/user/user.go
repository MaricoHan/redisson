package user

import (
	"errors"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gorm.io/gorm"
)

type IUserRepo interface {
	CreateUserTable() error
	Insert(user entity.User) error
	QueryUserById(userId uint64) (user *entity.User, err error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUserTable() error {
	if !u.db.Migrator().HasTable(&entity.User{}) {
		err := u.db.AutoMigrate(&entity.User{})
		if err != nil {
			return err
		}
		if u.db.Migrator().HasTable(&entity.User{}) {
			return nil
		} else {
			return errors.New("create user table failed")
		}
	} else {
		return nil
	}
}

func (u *UserRepo) Insert(user entity.User) error {
	return u.db.Create(&user).Error
}

func (u *UserRepo) QueryChainById(chainId uint64) (chain *entity.Chain, err error) {
	err = u.db.Where("id = ?", chainId).Find(&chain).Error
	return chain, err

}
