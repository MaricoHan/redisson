package handlers

import (
	"context"
	"strings"
	"time"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
)

type IAccount interface {
	CreateAccount(ctx context.Context, _ interface{}) (interface{}, error)
	Accounts(ctx context.Context, _ interface{}) (interface{}, error)
	AccountsHistory(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewAccount(svc *service.Account) IAccount {
	return newAccount(svc)
}

type account struct {
	base
	pageBasic
	svc *service.Account
}

func newAccount(svc *service.Account) *account {
	return &account{svc: svc}
}

// CreateAccount Create one or more accounts
// return creation result
func (h account) CreateAccount(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateAccountRequest)
	params := dto.CreateAccountP{
		AppID: h.AppID(ctx),
		Count: req.Count,
	}
	if params.Count == 0 {
		params.Count = 1
	}
	if params.Count > 1000 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "Invalid Count")
	}
	// 校验参数 end
	return h.svc.CreateAccount(params)
}

// Accounts return account list
func (h account) Accounts(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	params := dto.AccountsP{
		AppID:   h.AppID(ctx),
		Account: h.Account(ctx),
	}

	offset, err := h.Offset(ctx)
	if err != nil {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "Invalid Offset")
	}
	params.Offset = offset

	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "Invalid Limit")
	}
	params.Limit = limit

	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit > 50 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "Invalid Limit")
	}

	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		startDateTime, err := time.Parse(timeLayout, startDateR+" 00:00:00")
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "Invalid StartDate")
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, endDateR+" 23:59:59")
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "Invalid EndDate")
		}
		params.EndDate = &endDateTime
	}

	if params.EndDate != nil && params.StartDate != nil {
		if !params.EndDate.After(*params.StartDate) {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "EndDate before StartDate")
		}
	}
	switch h.SortBy(ctx) {
	case "DATE_ASC":
		params.SortBy = "DATE_ASC"
	case "DATE_DESC":
		params.SortBy = "DATE_DESC"
	default:
		return nil, types.ErrParams
	}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.Accounts(params)
}

func (h account) Account(ctx context.Context) string {
	accountR := ctx.Value("account")
	if accountR == nil || accountR == "" {
		return ""
	}
	return accountR.(string)
}

func (h account) AccountsHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	params := dto.AccountsP{
		AppID:   h.AppID(ctx),
		Account: h.Account(ctx),
	}

	offset, err := h.Offset(ctx)
	if err != nil {
		return nil, types.ErrParams
	}
	params.Offset = offset

	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, types.ErrParams
	}
	params.Limit = limit

	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit > 50 {
		return nil, types.ErrParams
	}

	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		startDateTime, err := time.Parse(timeLayout, startDateR+" 00:00:00")
		if err != nil {
			return nil, types.ErrParams
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, endDateR+" 23:59:59")
		if err != nil {
			return nil, types.ErrParams
		}
		params.EndDate = &endDateTime
	}

	if params.EndDate != nil && params.StartDate != nil {
		if !params.EndDate.After(*params.StartDate) {
			return nil, types.ErrParams
		}
	}
	switch h.SortBy(ctx) {
	case "DATE_ASC":
		params.SortBy = "DATE_ASC"
	case "DATE_DESC":
		params.SortBy = "DATE_DESC"
	default:
		return nil, types.ErrParams
	}

	params.Module = h.module(ctx)
	params.Operation = h.operation(ctx)
	if params.Module != "" && params.Operation != "" {
		if params.Module == "account" && params.Operation != "add_gas" {
			return nil, types.ErrParams
		} else if params.Module == "nft" && !strings.Contains("transfer_class/mint/edit/transfer/burn", params.Operation) {
			return nil, types.ErrParams
		}
	}

	if params.Module == "" && params.Operation != "" {
		params.Operation = ""
	}
	return h.svc.AccountsHistory(params)
}

func (h account) module(ctx context.Context) string {
	module := ctx.Value("module")
	if module == nil || module == "" {
		return ""
	}
	return module.(string)
}

func (h account) operation(ctx context.Context) string {
	operation := ctx.Value("operation")
	if operation == nil || operation == "" {
		return ""
	}
	return operation.(string)
}
