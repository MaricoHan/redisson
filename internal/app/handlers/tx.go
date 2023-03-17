package handlers

import (
	"context"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
)

type ITx interface {
	TxResult(ctx context.Context, _ interface{}) (interface{}, error)
	TxQueueInfo(ctx context.Context, _ interface{}) (interface{}, error)
}

type Tx struct {
	base
	svc services.ITx
}

func NewTx(svc services.ITx) *Tx {
	return &Tx{svc: svc}
}

func (h *Tx) TxResult(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.TxResultByTxHash{
		OperationId: h.OperationId(ctx),
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
		AccessMode:  authData.AccessMode,
	}
	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.TxResult(ctx, params)
}

func (h *Tx) TxQueueInfo(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.TxQueueInfo{
		OperationId: h.OperationId(ctx),
		ProjectID:   authData.ProjectId,
		Module:      authData.Module,
		Code:        authData.Code,
	}
	// 校验参数 end
	return h.svc.TxQueueInfo(ctx, params)
}

func (h *Tx) OperationId(ctx context.Context) string {
	operationId := ctx.Value("operation_id")
	if operationId == nil {
		return ""
	}
	return operationId.(string)
}
