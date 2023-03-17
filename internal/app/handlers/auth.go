package handlers

import (
	"context"
	"fmt"

	"gitlab.bianjie.ai/avata/utils/errors"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
)

const (
	AUTHTYPEID    = "1"
	AUTHTYPEPHONE = "2"
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
	hashType := a.hashType(ctx)
	if hash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "hash"))
	}
	if projectID == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "project_id"))
	}
	if hashType == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "type"))
	}
	if hashType != AUTHTYPEID && hashType != AUTHTYPEPHONE {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(constant.ErrInvalidValue, "type"))
	}
	return a.svc.Verify(ctx, &vo.AuthVerify{
		Hash:      hash,
		ProjectID: projectID,
		Type:      hashType,
	})
}

// GetUser 身份信息查询
func (a *Auth) GetUser(ctx context.Context, request interface{}) (interface{}, error) {
	hash := a.hash(ctx)
	projectID := a.projectID(ctx)
	hashType := a.hashType(ctx)
	if hash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "hash"))
	}
	if projectID == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "project_id"))
	}
	if hashType == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "type"))
	}
	if hashType != AUTHTYPEID && hashType != AUTHTYPEPHONE {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(constant.ErrInvalidValue, "type"))
	}
	return a.svc.GetUser(ctx, &vo.AuthGetUser{
		Hash:      hash,
		ProjectID: projectID,
		Type:      hashType,
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

func (a *Auth) hashType(ctx context.Context) string {
	hashType := ctx.Value("type")
	if hashType == nil {
		return ""
	}
	return hashType.(string)
}
