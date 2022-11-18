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
	hash := a.hash(ctx)
	projectID := a.projectID(ctx)
	if hash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "hash"))
	}
	if projectID == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "project_id"))
	}
	return a.svc.Verify(ctx, &vo.AuthVerify{
		Hash:      hash,
		ProjectID: projectID,
	})
}

// GetUser 身份信息查询
func (a *Auth) GetUser(ctx context.Context, request interface{}) (interface{}, error) {
	hash := a.hash(ctx)
	projectID := a.projectID(ctx)
	phoneHash := a.phoneHash(ctx)
	if hash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "hash"))
	}
	if projectID == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "project_id"))
	}
	if phoneHash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "phone_hash"))
	}
	return a.svc.GetUser(ctx, &vo.AuthGetUser{
		Hash:      hash,
		ProjectID: projectID,
		PhoneHash: phoneHash,
	})
}

func (a *Auth) hash(ctx context.Context) string {
	hash := ctx.Value("hash")
	if hash == nil {
		return ""
	}
	return hash.(string)
}

func (a *Auth) projectID(ctx context.Context) string {
	projectID := ctx.Value("project_id")
	if projectID == nil {
		return ""
	}
	return projectID.(string)
}

func (a *Auth) phoneHash(ctx context.Context) string {
	phoneHash := ctx.Value("phone_hash")
	if phoneHash == nil {
		return ""
	}
	return phoneHash.(string)
}
