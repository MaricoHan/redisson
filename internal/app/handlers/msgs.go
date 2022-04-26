package handlers

import (
	"context"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
)

type IMsgs interface {
	GetNFTHistory(ctx context.Context, _ interface{}) (interface{}, error)
	GetAccountHistory(ctx context.Context, _ interface{}) (interface{}, error)
}

type Msgs struct {
	base
	pageBasic
	svc services.IMsgs
}

func NewMsgs(svc services.IMsgs) *Msgs {
	return &Msgs{svc: svc}
}

func (h *Msgs) GetNFTHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)
	params := dto.NftOperationHistoryByNftId{
		ClassID:    h.ClassId(ctx),
		NftId:      h.NftId(ctx),
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Code:       authData.Code,
	}

	offset, err := h.Offset(ctx)
	if err != nil {
		return nil, err
	}
	params.Offset = offset

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

	params.Operation = h.Operation(ctx)

	return h.svc.GetNFTHistory(params)
}

func (h *Msgs) GetAccountHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.AccountsInfo{
		ChainID:         authData.ChainId,
		ProjectID:       authData.ProjectId,
		PlatFormID:      authData.PlatformId,
		Account:         h.Account(ctx),
		Module:          authData.Module,
		Code:            authData.Code,
		OperationModule: h.operationModule(ctx),
	}

	offset, err := h.Offset(ctx)
	if err != nil {
		return nil, err
	}
	params.Offset = offset

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
	params.Operation = h.operation(ctx)

	return h.svc.GetAccountHistory(params)
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

func (h *Msgs) Operation(ctx context.Context) string {
	operation := ctx.Value("operation")
	if operation == nil || operation == "" {
		return ""
	}
	return operation.(string)
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

func (h *Msgs) operationModule(ctx context.Context) string {
	module := ctx.Value("module")
	if module == nil || module == "" {
		return ""
	}
	return module.(string)
}

func (h *Msgs) operation(ctx context.Context) string {
	operation := ctx.Value("operation")
	if operation == nil || operation == "" {
		return ""
	}
	return operation.(string)
}
