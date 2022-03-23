package handlers

import (
	"context"
	"fmt"
	"time"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

type account struct {
	base
	pageBasic
	svc map[string]service.AccountService
}

type IAccount interface {
	Create(ctx context.Context, _ interface{}) (interface{}, error)
	Accounts(ctx context.Context, _ interface{}) (interface{}, error)
	AccountsHistory(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewAccount(svc ...*service.AccountBase) IAccount {
	return newAccountModule(svc)
}

func newAccountModule(svc []*service.AccountBase) *account {
	modules := make(map[string]service.AccountService, len(svc))
	for _, v := range svc {
		modules[v.Module] = v.Service
	}
	return &account{
		svc: modules,
	}
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

func (h account) Create(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateAccountRequest)

	authData := h.AuthData(ctx)
	params := dto.CreateAccountP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Count:      req.Count,
	}
	if params.Count == 0 {
		params.Count = 1
	}
	if params.Count < 1 || params.Count > 1000 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrCountLen)
	}
	//if config.Get().Server.Env == "prod" && params.Count > 10 {
	//	log.Error("create account", "params error:", "config.Get().Server.Env == \"prod\" && params.Count > 10")
	//	return nil, types.ErrParams
	//}
	// 校验参数 end
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Create(params)
}

func (h account) Accounts(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.AccountsP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Account:    h.Account(ctx),
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
		startDateTime, err := time.Parse(timeLayout, fmt.Sprintf("%s 00:00:00", startDateR))
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrStartDate)
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, fmt.Sprintf("%s 23:59:59", endDateR))
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrEndDate)
		}
		params.EndDate = &endDateTime
	}

	if params.EndDate != nil && params.StartDate != nil {
		if params.EndDate.Before(*params.StartDate) {
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
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Show(params)
}

func (h account) AccountsHistory(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.AccountsP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Account:    h.Account(ctx),
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
		startDateTime, err := time.Parse(timeLayout, fmt.Sprintf("%s 00:00:00", startDateR))
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrStartDate)
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, fmt.Sprintf("%s 23:59:59", endDateR))
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrEndDate)
		}
		params.EndDate = &endDateTime
	}

	if params.EndDate != nil && params.StartDate != nil {
		if params.EndDate.Before(*params.StartDate) {
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

	if params.Module == "" {
		params.Operation = ""
	} else {
		_, ok := ModuleOperation[params.Module]
		if !ok {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrModule)
		}
	}

	if params.Module != "" && params.Operation != "" {
		operation, ok := ModuleOperation[params.Module]
		if !ok {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrModule)
		}
		if _, ok = operation[params.Operation]; !ok {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOperation)
		}
	}
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.History(params)
}

func (h account) Account(ctx context.Context) string {
	accountR := ctx.Value("account")
	if accountR == nil || accountR == "" {
		return ""
	}
	return accountR.(string)
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
