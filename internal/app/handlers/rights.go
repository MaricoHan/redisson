package handlers

import (
	"context"
	service "gitlab.bianjie.ai/avata/open-api/internal/app/services"
)

type IRights interface {
	Register(ctx context.Context, request interface{}) (response interface{}, err error)
	EditRegister(ctx context.Context, request interface{}) (response interface{}, err error)
	QueryRegister(ctx context.Context, request interface{}) (response interface{}, err error)
}

type Rights struct {
	base
	svc service.IRights
}

func NewRights(svc service.IRights) *Rights {
	return &Rights{svc: svc}
}

func (r Rights) Register(ctx context.Context, request interface{}) (response interface{}, err error) {
	panic("implement me")
}

func (r Rights) EditRegister(ctx context.Context, request interface{}) (response interface{}, err error) {
	panic("implement me")
}

func (r Rights) QueryRegister(ctx context.Context, request interface{}) (response interface{}, err error) {
	panic("implement me")
}
