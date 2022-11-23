package service_redirect_url

import (
	"gorm.io/gorm"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
)

type IServiceRedirectUrlRepo interface {
	GetServiceRedirectUrlByProjectID(projectID uint64) (user entity.ServiceRedirectUrl, err error)
}

type ServiceRedirectUrlRepo struct {
	db *gorm.DB
}

func NewServiceRedirectUrlRepo(db *gorm.DB) *ServiceRedirectUrlRepo {
	return &ServiceRedirectUrlRepo{db: db}
}

func (p *ServiceRedirectUrlRepo) GetServiceRedirectUrlByProjectID(projectID uint64) (sru entity.ServiceRedirectUrl, err error) {
	err = p.db.Select("id,project_id, url").Where("project_id = ?", projectID).Find(&sru).Error
	return sru, err
}
