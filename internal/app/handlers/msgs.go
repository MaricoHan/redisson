package handlers

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.bianjie.ai/avata/utils/errors"
	"gitlab.bianjie.ai/avata/utils/errors/common"

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
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}
	nftId, err := h.NftId(ctx)
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

	operation, err := h.operation(ctx)
	if err != nil {
		return nil, err
	}
	params.Operation = uint32(operation)

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

	module, err := h.operationModule(ctx)
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
	operation, err := h.operation(ctx)
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

func (h *Msgs) ClassId(ctx context.Context) string {
	classId := ctx.Value("class_id")

	if classId == nil {
		return ""
	}
	return classId.(string)

}

func (h *Msgs) NftId(ctx context.Context) (uint64, error) {
	v := ctx.Value("nft_id")
	if v == nil {
		return 0, errors.New(errors.NotFound, "")
	}
	res, err := strconv.ParseUint(v.(string), 10, 64)
	if err != nil {
		return 0, errors.New(errors.NotFound, fmt.Sprintf("%s, nft_id: %s not found", errors.ErrResourceNotFound, v.(string)))
	}

	return res, nil
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

func (h *Msgs) operationModule(ctx context.Context) (uint32, error) {
	v := ctx.Value("module")
	if v == nil {
		return 0, nil
	}
	m := v.(string)

	res, err := strconv.ParseUint(m, 10, 64)
	if err != nil {
		return 0, errors.ErrModules
	}

	return uint32(res), nil
}

func (h *Msgs) operation(ctx context.Context) (uint32, error) {
	v := ctx.Value("operation")
	if v == nil {
		return 0, nil
	}

	res, err := strconv.ParseUint(v.(string), 10, 64)
	if err != nil {
		return 0, errors.New(errors.ClientParams, errors.ErrOperation)
	}

	return uint32(res), nil
}
