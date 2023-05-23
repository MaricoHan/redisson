package evm

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/evm"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/evm"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/evm"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IContract interface {
	CreateCall(ctx context.Context, _ interface{}) (interface{}, error)
	ShowCall(ctx context.Context, _ interface{}) (interface{}, error)
}

type Contract struct {
	base.Base
	base.PageBasic
	svc evm.IContract
}

func NewContract(svc evm.IContract) *Contract {
	return &Contract{svc: svc}
}

func (h *Contract) CreateCall(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.CreateContractCallRequest)
	from := strings.TrimSpace(req.From)
	to := strings.TrimSpace(req.To)
	data := strings.TrimSpace(req.Data)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.CreateContractCall{
		OperationId: operationId,
		ProjectID:   authData.ProjectId,
		Module:      authData.Module,
		Code:        authData.Code,
		AccessMode:  authData.AccessMode,
		From:        from,
		To:          to,
		Data:        data,
		GasLimit:    req.GasLimit,
		Estimation:  req.Estimation,
	}
	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.CreateCall(ctx, params)
}

func (h *Contract) ShowCall(ctx context.Context, request interface{}) (interface{}, error) {
	to := strings.TrimSpace(h.To(ctx))
	data := strings.TrimSpace(h.Data(ctx))

	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.ShowContractCall{
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
		To:         to,
		Data:       data,
	}
	return h.svc.ShowCall(ctx, params)
}

func (h *Contract) From(ctx context.Context) string {
	from := ctx.Value("from")
	if from == nil {
		return ""
	}
	return from.(string)
}

func (h *Contract) To(ctx context.Context) string {
	to := ctx.Value("to")
	if to == nil {
		return ""
	}
	return to.(string)
}

func (h *Contract) Data(ctx context.Context) string {
	data := ctx.Value("data")
	if data == nil {
		return ""
	}
	return data.(string)
}
