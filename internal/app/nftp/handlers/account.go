package handlers

import (
	"context"
	"time"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
)

type IAccount interface {
	CreateAccount(ctx context.Context, _ interface{}) (interface{}, error)
	Accounts(ctx context.Context, _ interface{}) (interface{}, error)
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
	req := request.(vo.CreateAccountRequest)
	params := dto.CreateAccountP{
		AppID: h.AppID(ctx),
		Count: req.Count,
	}
	if params.Count == 0 {
		params.Count = 1
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
		PageP: dto.PageP{
			Offset: h.Offset(ctx),
			Limit:  h.Limit(ctx),
		},
	}
	if params.Offset == 0 {
		params.Offset = 1
	}

	if params.Limit == 0 {
		params.Limit = 10
	}
	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		startDateTime, err := time.Parse(timeLayout, startDateR)
		if err != nil {
			return nil, types.ErrParams
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, endDateR)
		if err != nil {
			return nil, types.ErrParams
		}
		params.EndDate = &endDateTime
	}

	if params.EndDate != nil && params.StartDate != nil {
		if params.EndDate.After(*params.StartDate) {
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

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.Accounts(params)
}

func (h account) Account(ctx context.Context) string {
	accountR := ctx.Value("account")
	if accountR == nil {
		return ""
	}
	return accountR.(string)
}
