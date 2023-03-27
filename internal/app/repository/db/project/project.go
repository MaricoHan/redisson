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
	err = p.db.Model(&entity.Project{}).
		Omit(
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

func (p *ProjectRepo) ExistServices(projectId, serviceType uint) (bool, error) {
	var Ids []uint64
	var projectServices []*entity.ProjectServices
	if err := p.db.Debug().Model(&entity.ProjectServices{}).Select("service_id").Where("project_id = ?", projectId).Find(&projectServices).Error; err != nil {
		return false, err
	}

	for _, v := range projectServices {
		Ids = append(Ids, v.ServiceId)
	}

	var services []*entity.Services
	if err := p.db.Debug().Model(&entity.Service{}).Where("id IN ? AND type = ?", Ids, serviceType).Find(&services).Error; err != nil {
		return false, err
	}

	if len(services) > 0 {
		return true, nil
	}
	return false, nil
}
