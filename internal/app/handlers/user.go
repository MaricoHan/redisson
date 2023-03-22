package handlers

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
)

type IUser interface {
	CreateUsers(ctx context.Context, request interface{}) (interface{}, error)
	UpdateUsers(ctx context.Context, request interface{}) (interface{}, error)
}

type User struct {
	base
	pageBasic
	svc services.IUser
}

func NewUser(svc services.IUser) *User {
	return &User{svc: svc}
}

func (u User) CreateUsers(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.CreateUserRequest)

	// 校验参数
	req.Individual.Name = strings.TrimSpace(req.Individual.Name)

	authData := u.AuthData(ctx)
	params := dto.CreateUsers{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}
	return u.svc.CreateUsers(ctx, params)
}

func (u User) UpdateUsers(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.UpdateUserRequest)

	// 校验参数
	userId := strings.TrimSpace(req.UserId)

	authData := u.AuthData(ctx)
	params := dto.UpdateUsers{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
		UserId:     userId,
	}
	return u.svc.UpdateUsers(ctx, params)
}
