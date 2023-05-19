package evm

import (
	"context"
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	evm2 "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/evm"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/evm"
	"gitlab.bianjie.ai/avata/utils/errors"
	"gitlab.bianjie.ai/avata/utils/errors/common"
)

type IMsgs interface {
	GetNFTHistory(ctx context.Context, _ interface{}) (interface{}, error)
	GetAccountHistory(ctx context.Context, _ interface{}) (interface{}, error)
}

type Msgs struct {
	handlers.Base
	handlers.PageBasic
	NFT
	svc evm.IMsgs
}

func NewMsgs(svc evm.IMsgs) *Msgs {
	return &Msgs{svc: svc}
}

func (h *Msgs) GetNFTHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)
	params := evm2.NftOperationHistoryByNftId{
		ClassID:    h.ClassID(ctx),
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}
	nftId, err := h.NftID(ctx)
	if err != nil {
		return nil, err
	}

	params.NftId = nftId

	params.PageKey = h.PageKey(ctx)
	countTotal, err := h.CountTotal(ctx)
	if err != nil {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(common.ERR_INVALID_VALUE, "count_total"))
	}
	params.CountTotal = countTotal

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
	params.TxHash = h.Txhash(ctx)

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
	countTotal, err := h.CountTotal(ctx)
	if err != nil {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(common.ERR_INVALID_VALUE, "count_total"))
	}
	params.CountTotal = countTotal
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
	operation, err := h.Operation(ctx)
	if err != nil {
		return nil, err
	}
	if module > 0 {
		params.Operation = operation
	}
	params.TxHash = h.Txhash(ctx)

	return h.svc.GetAccountHistory(ctx, params)
}

func (h *Msgs) MTId(ctx context.Context) string {
	nftId := ctx.Value("mt_id")
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
