package project

import (
	"gorm.io/gorm"

	constant "gitlab.bianjie.ai/avata/open-api/internal/app/models"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
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
	err = p.db.Omit(
		entity.ProjectFields.Code,
		entity.ProjectFields.Name,
		entity.ProjectFields.Description,
		entity.ProjectFields.CreatedAt,
		entity.ProjectFields.UpdatedAt,
	).Where(entity.ProjectFields.ApiKey, apiKey).
		Where(entity.ProjectFields.Status, constant.ProjectStatusEnable).
		Find(&project).Error

	return
}

func (p *ProjectRepo) GetProjectByCode(code string) (project entity.Project, err error) {
	err = p.db.Omit(
		entity.ProjectFields.Code,
		entity.ProjectFields.Name,
		entity.ProjectFields.Description,
		entity.ProjectFields.CreatedAt,
		entity.ProjectFields.UpdatedAt,
	).Where(entity.ProjectFields.Code, code).
		Where(entity.ProjectFields.Status, constant.ProjectStatusEnable).
		Find(&project).Error
	return
}
