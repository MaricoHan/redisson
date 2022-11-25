package handlers

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IMTClass interface {
	CreateMTClass(ctx context.Context, _ interface{}) (interface{}, error)
	TransferMTClass(ctx context.Context, _ interface{}) (interface{}, error)
	List(ctx context.Context, _ interface{}) (interface{}, error)
	Show(ctx context.Context, _ interface{}) (interface{}, error)
}

type MTClass struct {
	base
	pageBasic
	svc services.IMTClass
}

func NewMTClass(svc services.IMTClass) *MTClass {
	return &MTClass{svc: svc}
}

func (h MTClass) CreateMTClass(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateMTClassRequest)

	name := strings.TrimSpace(req.Name)
	data := strings.TrimSpace(req.Data)
	owner := strings.TrimSpace(req.Owner)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	if len([]rune(owner)) > 128 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOwnerLen)
	}

	authData := h.AuthData(ctx)
	params := dto.CreateMTClass{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Name:        name,
		Data:        data,
		Owner:       owner,
		Tag:         tagBytes,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	return h.svc.CreateMTClass(ctx, params)
}

func (h MTClass) TransferMTClass(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferMTClassRequest)
	recipient := strings.TrimSpace(req.Recipient)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if recipient == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrRecipient)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}

	// 校验参数 end
	authData := h.AuthData(ctx)
	params := dto.TransferMTClass{
		ClassID:     h.Id(ctx),
		Owner:       h.Owner(ctx),
		Recipient:   recipient,
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Tag:         tagBytes,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	return h.svc.TransferMTClass(ctx, params)
}

// List return class list
func (h MTClass) List(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.MTClassListRequest{
		ClassId:    h.Id(ctx),
		ClassName:  h.Name(ctx),
		Owner:      h.Owner(ctx),
		TxHash:     h.TxHash(ctx),
		ProjectID:  authData.ProjectId,
		ChainID:    authData.ChainId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
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

	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		params.StartDate = startDateR
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		params.EndDate = endDateR
	}

	params.SortBy = h.SortBy(ctx)

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.List(ctx, &params)
}

// Show return class
func (h MTClass) Show(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.MTClassShowRequest{
		ProjectID:  authData.ProjectId,
		ClassID:    h.Id(ctx),
		Status:     h.Status(ctx),
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.Show(ctx, &params)
}

func (h MTClass) Id(ctx context.Context) string {
	val := ctx.Value("id")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MTClass) Name(ctx context.Context) string {
	val := ctx.Value("name")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MTClass) Owner(ctx context.Context) string {
	val := ctx.Value("owner")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MTClass) TxHash(ctx context.Context) string {
	val := ctx.Value("tx_hash")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MTClass) Timestamp(ctx context.Context) string {
	val := ctx.Value("timestamp")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MTClass) Status(ctx context.Context) string {
	val := ctx.Value("status")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MTClass) MtCount(ctx context.Context) uint64 {
	val := ctx.Value("mt_count")
	if val == nil {
		return 0
	}
	return val.(uint64)
}

func (h MTClass) ClassID(ctx context.Context) string {
	classId := ctx.Value("mt_class_id")
	if classId == nil {
		return ""
	}
	return classId.(string)
}
