package l2

import (
	"context"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/l2"
	service "gitlab.bianjie.ai/avata/open-api/internal/app/services/l2"
)

type IDict interface {
	ListTxTypes(ctx context.Context, request interface{}) (interface{}, error)
}

type Dict struct {
	handlers.Base
	svc service.IDict
}

func NewDict(svc service.IDict) IDict {
	return Dict{svc: svc}
}

var _ IDict = Dict{}

func (d Dict) ListTxTypes(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := d.AuthData(ctx)
	params := &dto.ListTxTypes{
		Code:   authData.Code,
		Module: authData.Module,
	}
	return d.svc.ListTxTypes(ctx, params)
}
