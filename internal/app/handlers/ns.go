package handlers

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INs interface {
	CreateDomain(ctx context.Context, _ interface{}) (interface{}, error)
	Domains(ctx context.Context, _ interface{}) (interface{}, error)
}

type Ns struct {
	base
	svc services.INs
}

func NewNs(svc services.INs) *Ns {
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
		OperationId: h.OperationId(ctx),
		ProjectID:   authData.ProjectId,
		Module:      authData.Module,
		Code:        authData.Code,
		AccessMode:  authData.AccessMode,
		Name:        name,
		Owner:       owner,
	}
	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.CreateDomain(ctx, params)
}

func (h *Ns) Domains(ctx context.Context, request interface{}) (interface{}, error) {
	name := strings.TrimSpace(h.Name(ctx))
	tld := strings.TrimSpace(h.Tld(ctx))

	if name == "" {
		// todo
		return nil, errors2.New(errors2.ClientParams, "empty name")
	}
	if tld == "" {
		// todo
		return nil, errors2.New(errors2.ClientParams, "empty tld")
	}

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
	// 校验参数 end
	return h.svc.Domains(ctx, params)
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
