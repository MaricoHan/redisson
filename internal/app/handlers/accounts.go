package handlers

import (
	"context"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	"gitlab.bianjie.ai/avata/utils/errors"
)

type IAccount interface {
	CreateAccount(ctx context.Context, _ interface{}) (interface{}, error)
	GetAccounts(ctx context.Context, _ interface{}) (interface{}, error)
}

type Account struct {
	base
	pageBasic
	svc services.IAccount
}

func NewAccount(svc services.IAccount) *Account {
	return &Account{svc: svc}
}

func (h *Account) CreateAccount(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateAccountRequest)

	authData := h.AuthData(ctx)
	params := dto.CreateAccount{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Count:      req.Count,
		Module:     authData.Module,
		Code:       authData.Code,
	}
	if params.Count < 0 || params.Count > 1000 {
		return nil, errors.New(errors.ClientParams, errors.ErrCountLen)
	}
	if params.Count == 0 {
		params.Count = 1
	}
	return h.svc.CreateAccount(params)
}

func (h *Account) GetAccounts(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.AccountsInfo{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Account:    h.Account(ctx),
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

	return h.svc.GetAccounts(params)
}

func (h *Account) Account(ctx context.Context) string {
	accountR := ctx.Value("account")
	if accountR == nil || accountR == "" {
		return ""
	}
	return accountR.(string)
}
