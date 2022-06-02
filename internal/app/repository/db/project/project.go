package project

import (
	"errors"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gorm.io/gorm"
)

type IProjectRepo interface {
	CreateChainTable() error
	GetProjectByApiKey(apiKey string) (project *entity.Project, err error)
	Insert(project entity.Project) error
	//GetProjectByUserId(userId int64)(projects []*entity.Project, err error)
}

type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (p *ProjectRepo) CreateProjectTable() error {
	if !p.db.Migrator().HasTable(&entity.Project{}) {
		err := p.db.AutoMigrate(&entity.Project{})
		if err != nil {
			return err
		}
		if p.db.Migrator().HasTable(&entity.Project{}) {
			return nil
		} else {
			return errors.New("create project table failed")
		}
	} else {
		return nil
	}
}

func (p *ProjectRepo) Insert(project entity.Project) error {
	return p.db.Create(&project).Error
}

func (p *ProjectRepo) GetProjectByApiKey(apiKey string) (project entity.Project, err error) {
	err = p.db.Where("api_key=?", apiKey).Find(&project).Error
	return project, err
}
