package handlers

import (
	"context"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	"gitlab.bianjie.ai/avata/open-api/utils"
	"gitlab.bianjie.ai/avata/utils/errors"
	"strings"
)

type IAccount interface {
	BatchCreateAccount(ctx context.Context, _ interface{}) (interface{}, error)
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

// BatchCreateAccount 批量创建链账户
func (h *Account) BatchCreateAccount(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.BatchCreateAccountRequest)

	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors.New(errors.ClientParams, errors.ErrOperationID)
	}

	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors.New(errors.ClientParams, errors.ErrOperationIDLen)
	}

	authData := h.AuthData(ctx)
	params := dto.BatchCreateAccount{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Count:       req.Count,
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: operationId,
	}
	if params.Count < 0 || params.Count > 1000 {
		return nil, errors.New(errors.ClientParams, errors.ErrCountLen)
	}
	if params.Count == 0 {
		params.Count = 1
	}
	return h.svc.BatchCreateAccount(params)
}

// CreateAccount 单个创建链账户
func (h *Account) CreateAccount(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateAccountRequest)
	name := strings.TrimSpace(req.Name)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors.New(errors.ClientParams, errors.ErrOperationID)
	}
	if name == "" {
		return nil, errors.New(errors.ClientParams, errors.ErrName)
	}

	if len([]rune(name)) < 1 || len([]rune(name)) > 20 {
		return nil, errors.New(errors.ClientParams, errors.ErrAccountNameLen)
	}
	if !utils.StrNameCheck(name) {
		return nil, errors.New(errors.ClientParams, errors.ErrNameFormat)
	}

	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors.New(errors.ClientParams, errors.ErrOperationIDLen)
	}
	authData := h.AuthData(ctx)
	params := dto.CreateAccount{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Name:        name,
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: operationId,
	}

	return h.svc.CreateAccount(params)
}

func (h *Account) GetAccounts(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.AccountsInfo{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Account:     h.Account(ctx),
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: h.OperationID(ctx),
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

func (h *Account) OperationID(ctx context.Context) string {
	OperationID := ctx.Value("operation_id")
	if OperationID == nil || OperationID == "" {
		return ""
	}
	return OperationID.(string)
}