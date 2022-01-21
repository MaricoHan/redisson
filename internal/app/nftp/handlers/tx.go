package handlers

import (
	"context"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
)

type ITx interface {
	TxResultByTxHash(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewTx(svc *service.Tx) ITx {
	return newTx(svc)
}

type tx struct {
	base
	svc *service.Tx
}

func newTx(svc *service.Tx) *tx {
	return &tx{svc: svc}
}

// TxResultByTxHash query txresult by txhash
func (h tx) TxResultByTxHash(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	params := dto.TxResultByTxHashP{
		Hash:  h.Hash(ctx),
		AppID: h.AppID(ctx),
	}
	if params.Hash == "" {
		params.Hash = ""
	}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.TxResultByTxHash(params)
}

func (h tx) Hash(ctx context.Context) string {
	hash := ctx.Value("hash")
	if hash == nil {
		return ""
	}
	return hash.(string)
}
