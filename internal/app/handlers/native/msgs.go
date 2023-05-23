package native

import (
	"context"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/native"
)

type IMsgs interface {
	GetNFTHistory(ctx context.Context, _ interface{}) (interface{}, error)
	GetAccountHistory(ctx context.Context, _ interface{}) (interface{}, error)
	GetMTHistory(ctx context.Context, _ interface{}) (interface{}, error)
}

type Msgs struct {
	base.Base
	base.PageBasic
	svc native.IMsgs
}

func NewMsgs(svc native.IMsgs) *Msgs {
	return &Msgs{svc: svc}
}

func (h *Msgs) GetNFTHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)
	params := nft.NftOperationHistoryByNftId{
		ClassID:    h.ClassId(ctx),
		NftId:      h.NftId(ctx),
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}

	params.PageKey = h.PageKey(ctx)
	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, err
	}
	params.Limit = limit

	if params.Limit == 0 {
		params.Limit = 10
	}

	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		params.StartDate = startDateR
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		params.EndDate = endDateR
	}

	params.SortBy = h.SortBy(ctx)

	params.Signer = h.Signer(ctx)
	params.Txhash = h.Txhash(ctx)

	operation, err := h.Operation(ctx)
	if err != nil {
		return nil, err
	}
	params.Operation = operation

	return h.svc.GetNFTHistory(ctx, params)
}

func (h *Msgs) GetAccountHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.AccountsInfo{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Account:    h.Account(ctx),
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}

	module, err := h.OperationModule(ctx)
	if err != nil {
		return nil, err
	}
	params.OperationModule = module

	params.PageKey = h.PageKey(ctx)
	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, err
	}
	params.Limit = limit

	if params.Limit == 0 {
		params.Limit = 10
	}

	startDateR := h.StartDate(ctx)

	if startDateR != "" {
		params.StartDate = startDateR
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		params.EndDate = endDateR
	}

	params.SortBy = h.SortBy(ctx)

	params.TxHash = h.Txhash(ctx)

	operation, err := h.Operation(ctx)
	if err != nil {
		return nil, err
	}
	if params.OperationModule > 0 {
		params.Operation = operation
	}
	return h.svc.GetAccountHistory(ctx, params)
}

func (h *Msgs) GetMTHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)
	params := mt.MTOperationHistoryByMTId{
		ClassID:    h.ClassId(ctx),
		MTId:       h.MTId(ctx),
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}

	params.PageKey = h.PageKey(ctx)
	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, err
	}
	params.Limit = limit

	if params.Limit == 0 {
		params.Limit = 10
	}

	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		params.StartDate = startDateR
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		params.EndDate = endDateR
	}

	params.SortBy = h.SortBy(ctx)
	params.Signer = h.Signer(ctx)
	params.Txhash = h.Txhash(ctx)
	operation, err := h.Operation(ctx)
	if err != nil {
		return nil, err
	}
	params.Operation = operation

	return h.svc.GetMTHistory(ctx, params)
}

func (h *Msgs) MTId(ctx context.Context) string {
	nftId := ctx.Value("mt_id")
	if nftId == nil {
		return ""
	}
	return nftId.(string)
}

func (h *Msgs) ClassId(ctx context.Context) string {
	classId := ctx.Value("class_id")

	if classId == nil {
		return ""
	}
	return classId.(string)

}

func (h *Msgs) NftId(ctx context.Context) string {
	nftId := ctx.Value("nft_id")
	if nftId == nil {
		return ""
	}
	return nftId.(string)
}

func (h *Msgs) Signer(ctx context.Context) string {
	signer := ctx.Value("signer")
	if signer == nil || signer == "" {
		return ""
	}
	return signer.(string)
}

func (h *Msgs) Txhash(ctx context.Context) string {
	txhash := ctx.Value("tx_hash")
	if txhash == nil || txhash == "" {
		return ""
	}
	return txhash.(string)
}

func (h *Msgs) Account(ctx context.Context) string {
	accountR := ctx.Value("account")
	if accountR == nil || accountR == "" {
		return ""
	}
	return accountR.(string)
}
