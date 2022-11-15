package handlers

import (
	"context"
	"fmt"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/utils/errors"

	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
)

type IAuth interface {
	Verify(ctx context.Context, _ interface{}) (interface{}, error)
	GetUser(ctx context.Context, _ interface{}) (interface{}, error)
}

type Auth struct {
	base
	pageBasic
	svc services.IAuth
}

func NewAuth(svc services.IAuth) *Auth {
	return &Auth{svc: svc}
}

// Verify 身份信息验证
func (a *Auth) Verify(ctx context.Context, request interface{}) (interface{}, error) {
	params := request.(*vo.AuthVerify)
	if params.Hash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "hash"))
	}
	if params.ProjectID == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "project_id"))
	}
	return a.svc.Verify(ctx, params)
}

// GetUser 身份信息查询
func (a *Auth) GetUser(ctx context.Context, request interface{}) (interface{}, error) {
	params := request.(*vo.AuthGetUser)
	if params.Hash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "hash"))
	}
	if params.ProjectID == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "project_id"))
	}
	return a.svc.GetUser(ctx, params)
}
