package handlers

import (
	"context"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

type ITx interface {
	TxResultByTxHash(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewTx(svc ...*service.TXBase) ITx {
	return newTx(svc)
}

type tx struct {
	base
	svc map[string]service.TXService
}

func newTx(svc []*service.TXBase) *tx {
	modules := make(map[string]service.TXService, len(svc))
	for _, v := range svc {
		modules[v.Module] = v.Service
	}
	return &tx{
		svc: modules,
	}
}

// TxResultByTxHash query txresult by txhash
func (h tx) TxResultByTxHash(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.TxResultByTxHashP{
		TaskId:     h.TaskId(ctx),
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
	}
	// 校验参数 end
	// 业务数据入库的地方
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Show(params)
}

func (h tx) TaskId(ctx context.Context) string {
	taskid := ctx.Value("task_id")
	if taskid == nil {
		return ""
	}
	return taskid.(string)
}
