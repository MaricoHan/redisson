package handlers

import (
	"context"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
)

type IMTMsgs interface {
	GetMTHistory(ctx context.Context, _ interface{}) (interface{}, error)
}

type MTMsgs struct {
	base
	pageBasic
	svc services.IMTMsgs
}

func NewMTMsgs(svc services.IMTMsgs) *MTMsgs {
	return &MTMsgs{svc: svc}
}

func (h *MTMsgs) GetMTHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)
	params := dto.MTOperationHistoryByMTId{
		ClassID:    h.MTClassId(ctx),
		MTId:      h.MTId(ctx),
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

	return h.svc.GetMTHistory(params)
}

func (h *MTMsgs) MTClassId(ctx context.Context) string {
	classId := ctx.Value("mt_class_id")

	if classId == nil {
		return ""
	}
	return classId.(string)

}

func (h *MTMsgs) MTId(ctx context.Context) string {
	nftId := ctx.Value("mt_id")
	if nftId == nil {
		return ""
	}
	return nftId.(string)
}

func (h *MTMsgs) Signer(ctx context.Context) string {
	signer := ctx.Value("signer")
	if signer == nil || signer == "" {
		return ""
	}
	return signer.(string)
}

func (h *MTMsgs) Operation(ctx context.Context) string {
	operation := ctx.Value("operation")
	if operation == nil || operation == "" {
		return ""
	}
	return operation.(string)
}

func (h *MTMsgs) Txhash(ctx context.Context) string {
	txhash := ctx.Value("tx_hash")
	if txhash == nil || txhash == "" {
		return ""
	}
	return txhash.(string)
}
