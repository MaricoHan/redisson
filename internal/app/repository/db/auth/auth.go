package auth

import (
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gorm.io/gorm"
)

type IProjectAuthRepo interface {
	GetProjectPermission(pid int) (list []entity.Permission, err error)
}

type ProjectAuthRepo struct {
	db *gorm.DB
}

func NewProjectAuthRepo(db *gorm.DB) *ProjectAuthRepo {
	return &ProjectAuthRepo{db: db}
}

func (repo *ProjectAuthRepo) GetProjectPermission(pid int) (list []entity.Permission, err error) {
	// 查询项目拥有的服务id
	subquery := repo.db.Select(entity.ProjectXServiceFields.ServiceId).Where(entity.ProjectXServiceFields.ProjectId, pid).Table(constant.MysqlProjectXServicesTable)
	// 查询服务拥有的权限id
	subquery2 := repo.db.Select(entity.ServiceXPermissoinFields.PermissionId).Where(fmt.Sprintf("%s in (?)", entity.ServiceXPermissoinFields.ServiceId), subquery).Table(constant.MysqlServiceXPermissoinTable)
	// 查询权限id对应的权限详情
	err = repo.db.Omit(entity.PermissoinFields.CreatedAt, entity.PermissoinFields.UpdatedAt).Where(fmt.Sprintf("%s in (?)", entity.PermissoinFields.ID), subquery2).Table(constant.MysqlPermissoinTable).Order(entity.PermissoinFields.Priority).Find(&list).Error

	return list, err
}
