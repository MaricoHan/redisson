package evm

import (
	"context"
	"fmt"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/evm"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/evm"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/evm"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"
	"gitlab.bianjie.ai/avata/utils/errors/v2/common"
)

type INs interface {
	Domains(ctx context.Context, _ interface{}) (interface{}, error)
	UserDomains(ctx context.Context, _ interface{}) (interface{}, error)
	CreateDomain(ctx context.Context, _ interface{}) (interface{}, error)
	TransferDomain(ctx context.Context, _ interface{}) (interface{}, error)
}

type Ns struct {
	base.Base
	base.PageBasic
	svc evm.INs
}

func NewNs(svc evm.INs) *Ns {
	return &Ns{svc: svc}
}

func (h *Ns) CreateDomain(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.CreateDomainRequest)
	name := strings.TrimSpace(req.Name)
	owner := strings.TrimSpace(req.Owner)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.CreateDomain{
		OperationId: operationId,
		ProjectID:   authData.ProjectId,
		Module:      authData.Module,
		Code:        authData.Code,
		AccessMode:  authData.AccessMode,
		Name:        name,
		Owner:       owner,
		Duration:    req.Duration,
	}
	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.CreateDomain(ctx, params)
}

func (h *Ns) TransferDomain(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.TransferDomainRequest)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.TransferDomain{
		Name:        h.Name(ctx),
		Owner:       h.Owner(ctx),
		OperationId: operationId,
		ProjectID:   authData.ProjectId,
		Module:      authData.Module,
		Code:        authData.Code,
		AccessMode:  authData.AccessMode,
		Recipient:   req.Recipient,
	}
	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.TransferDomain(ctx, params)
}

func (h *Ns) Domains(ctx context.Context, request interface{}) (interface{}, error) {
	name := strings.TrimSpace(h.Name(ctx))
	tld := strings.TrimSpace(h.Tld(ctx))

	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.Domains{
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
		Name:       name,
		Tld:        tld,
	}
	return h.svc.Domains(ctx, params)
}

func (h *Ns) UserDomains(ctx context.Context, request interface{}) (interface{}, error) {
	name := strings.TrimSpace(h.Name(ctx))
	tld := strings.TrimSpace(h.Tld(ctx))
	owner := strings.TrimSpace(h.Owner(ctx))

	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.Domains{
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
		Owner:      owner,
		Name:       name,
		Tld:        tld,
	}
	params.PageKey = h.PageKey(ctx)
	countTotal, err := h.CountTotal(ctx)
	if err != nil {
		return nil, errors2.New(errors2.ClientParams, fmt.Sprintf(common.ERR_INVALID_VALUE, "count_total"))
	}
	params.CountTotal = countTotal
	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, err
	}
	params.Limit = limit
	// 校验参数 end
	return h.svc.UserDomains(ctx, params)
}

func (h *Ns) OperationId(ctx context.Context) string {
	operationId := ctx.Value("operation_id")
	if operationId == nil {
		return ""
	}
	return operationId.(string)
}

func (h *Ns) Name(ctx context.Context) string {
	name := ctx.Value("name")
	if name == nil {
		return ""
	}
	return name.(string)
}

func (h *Ns) Tld(ctx context.Context) string {
	tld := ctx.Value("tld")
	if tld == nil {
		return ""
	}
	return tld.(string)
}

func (h *Ns) Owner(ctx context.Context) string {
	owner := ctx.Value("owner")
	if owner == nil {
		return ""
	}
	return owner.(string)
}
