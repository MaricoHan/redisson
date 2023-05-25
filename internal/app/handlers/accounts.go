package handlers

import (
	"context"
	"fmt"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/utils/errors/v2"
	"gitlab.bianjie.ai/avata/utils/errors/v2/common"
)

type IAccount interface {
	BatchCreateAccount(ctx context.Context, _ interface{}) (interface{}, error)
	CreateAccount(ctx context.Context, _ interface{}) (interface{}, error)
	GetAccounts(ctx context.Context, _ interface{}) (interface{}, error)
}

type Account struct {
	base.Base
	base.PageBasic
	svc services.IAccount
}

func NewAccount(svc services.IAccount) *Account {
	return &Account{svc: svc}
}

// BatchCreateAccount 批量创建链账户
func (h *Account) BatchCreateAccount(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.BatchCreateAccountRequest)

	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors.New(errors.ClientParams, errors.ErrOperationID)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors.New(errors.ClientParams, errors.ErrOperationIDLen)
	}

	authData := h.AuthData(ctx)
	if authData.ExistWalletService {
		authData.Code = constant.Wallet
		authData.Module = constant.Server
	}
	params := dto.BatchCreateAccount{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Count:       uint32(req.Count),
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	if params.Count < 0 || params.Count > 1000 {
		return nil, errors.New(errors.ClientParams, errors.ErrCountLen)
	}
	if params.Count == 0 {
		params.Count = 1
	}
	return h.svc.BatchCreateAccount(ctx, params)
}

// CreateAccount 单个创建链账户
func (h *Account) CreateAccount(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateAccountRequest)

	name := strings.TrimSpace(req.Name)
	operationId := strings.TrimSpace(req.OperationID)
	userId := strings.TrimSpace(req.UserId)
	if operationId == "" {
		return nil, errors.New(errors.ClientParams, errors.ErrOperationID)
	}

	if name != "" && (len([]rune(name)) < 1 || len([]rune(name)) > 20) {
		return nil, errors.New(errors.ClientParams, errors.ErrAccountNameLen)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors.New(errors.ClientParams, errors.ErrOperationIDLen)
	}
	authData := h.AuthData(ctx)
	if authData.ExistWalletService {
		authData.Code = constant.Wallet
		authData.Module = constant.Server
	}
	params := dto.CreateAccount{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Name:        name,
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
		UserId:      userId,
	}

	return h.svc.CreateAccount(ctx, params)
}

func (h *Account) GetAccounts(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	if authData.ExistWalletService {
		authData.Code = constant.Wallet
		authData.Module = constant.Server
	}
	params := dto.AccountsInfo{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Account:     h.Account(ctx),
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: h.OperationID(ctx),
		Name:        h.Name(ctx),
		AccessMode:  authData.AccessMode,
		UserId:      h.UserId(ctx),
	}

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

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors.ErrNotImplemented
	}
	if fmt.Sprintf("%s-%s", params.Code, params.Module) == constant.WalletServer {
		return h.svc.GetUserAccounts(ctx, params)
	}
	return h.svc.GetAccounts(ctx, params)
}

func (h *Account) Account(ctx context.Context) string {
	accountR := ctx.Value("account")
	if accountR == nil || accountR == "" {
		return ""
	}
	return accountR.(string)
}

func (h *Account) OperationID(ctx context.Context) string {
	OperationID := ctx.Value("operation_id")
	if OperationID == nil || OperationID == "" {
		return ""
	}
	return OperationID.(string)
}

func (h *Account) Name(ctx context.Context) string {
	name := ctx.Value("name")
	if name == nil {
		return ""
	}
	return name.(string)
}

func (h *Account) UserId(ctx context.Context) string {
	accountR := ctx.Value("user_id")
	if accountR == nil || accountR == "" {
		return ""
	}
	return accountR.(string)
}
