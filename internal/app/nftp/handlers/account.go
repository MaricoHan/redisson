package handlers

import (
	"context"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"time"

	"gitlab.bianjie.ai/irita-paas/open-api/config"
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

var (
	// ModuleOperation 定义验证模块
	ModuleOperation = map[string]map[string]string{
		models.TMSGSModuleAccount: {
			models.TMSGSOperationAddGas: models.TMSGSOperationAddGas,
		},
		models.TMSGSModuleNFT: {
			models.TMSGSOperationIssueClass:    models.TMSGSOperationIssueClass,
			models.TMSGSOperationTransferClass: models.TMSGSOperationTransferClass,
			models.TMSGSOperationMint:          models.TMSGSOperationMint,
			models.TMSGSOperationEdit:          models.TMSGSOperationEdit,
			models.TMSGSOperationTransfer:      models.TMSGSOperationTransfer,
			models.TMSGSOperationBurn:          models.TMSGSOperationBurn,
		},
	}
)

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
	if params.Count < 1 || params.Count > 1000 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrCountLen)
	}
	if config.Get().Server.Env == "prod" && params.Count > 10 {
		return nil, types.ErrParams
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
		startDateTime, err := time.Parse(timeLayout, startDateR+" 00:00:00")
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrStartDate)
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, endDateR+" 23:59:59")
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrEndDate)
		}
		params.EndDate = &endDateTime
	}

	if params.EndDate != nil && params.StartDate != nil {
		if !params.EndDate.After(*params.StartDate) {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrDate)
		}
	}
	switch h.SortBy(ctx) {
	case "DATE_ASC":
		params.SortBy = "DATE_ASC"
	case "DATE_DESC":
		params.SortBy = "DATE_DESC"
	default:
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSortBy)
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
		startDateTime, err := time.Parse(timeLayout, startDateR+" 00:00:00")
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrStartDate)
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, endDateR+" 23:59:59")
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrEndDate)
		}
		params.EndDate = &endDateTime
	}

	if params.EndDate != nil && params.StartDate != nil {
		if !params.EndDate.After(*params.StartDate) {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrDate)
		}
	}
	switch h.SortBy(ctx) {
	case "DATE_ASC":
		params.SortBy = "DATE_ASC"
	case "DATE_DESC":
		params.SortBy = "DATE_DESC"
	default:
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSortBy)
	}

	params.Module = h.module(ctx)
	params.Operation = h.operation(ctx)
	if params.Module != "" && params.Operation != "" {
		operation, ok := ModuleOperation[params.Module]
		if !ok {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOperation)
		}
		if _, ok = operation[params.Operation]; !ok {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOperation)
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
