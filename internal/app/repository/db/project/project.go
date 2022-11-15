package project

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gorm.io/gorm"
)

type IProjectRepo interface {
	GetProjectByApiKey(apiKey string) (project entity.Project, err error)
	GetProjectByCode(code string) (project entity.Project, err error)
}

type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (p *ProjectRepo) GetProjectByApiKey(apiKey string) (project entity.Project, err error) {
	err = p.db.Select("id,chain_id,user_id,api_secret,api_key,access_mode").Where("api_key=?", apiKey).Find(&project).Error
	return project, err
}

func (p *ProjectRepo) GetProjectByCode(code string) (project entity.Project, err error) {
	err = p.db.Select("id,chain_id,user_id,api_secret,api_key,access_mode").Where("code=?", code).Find(&project).Error
	return project, err
}
