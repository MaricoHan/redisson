package l2

import (
	"context"
	"fmt"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/l2"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/l2"
	services "gitlab.bianjie.ai/avata/open-api/internal/app/services/l2"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"
	"gitlab.bianjie.ai/avata/utils/errors/v2/common"
)

type INftClass interface {
	CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error)
	Classes(ctx context.Context, _ interface{}) (interface{}, error)
	ClassByID(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftClassByID(ctx context.Context, request interface{}) (interface{}, error)
}

type NftClass struct {
	base.Base
	base.PageBasic
	svc services.INFTClass
}

func NewNFTClass(svc services.INFTClass) *NftClass {
	return &NftClass{svc: svc}
}

// CreateNftClass Create one nft class
// return creation result
func (h NftClass) CreateNftClass(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateNftClassRequest)

	name := strings.TrimSpace(req.Name)
	classId := strings.TrimSpace(req.ClassId)
	symbol := strings.TrimSpace(req.Symbol)
	description := strings.TrimSpace(req.Description)
	uri := strings.TrimSpace(req.Uri)
	uriHash := strings.TrimSpace(req.UriHash)
	data := strings.TrimSpace(req.Data)
	owner := strings.TrimSpace(req.Owner)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if name == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrName)
	}
	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if err := h.Base.UriCheck(uri); err != nil {
		return nil, err
	}

	if owner == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOwner)
	}

	if len([]rune(owner)) > 128 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOwnerLen)
	}

	authData := h.AuthData(ctx)
	params := dto.CreateNftClass{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Code:        authData.Code,
		Module:      authData.Module,
		AccessMode:  authData.AccessMode,
		Name:        name,
		ClassId:     classId,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Data:        data,
		Owner:       owner,
		OperationId: operationId,
	}
	return h.svc.Create(ctx, params)
}

// Classes return class list
func (h NftClass) Classes(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.NftClasses{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Id:         h.Id(ctx),
		Name:       h.Name(ctx),
		Owner:      h.Owner(ctx),
		TxHash:     h.TxHash(ctx),
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
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
	return h.svc.List(ctx, params)
}

// ClassByID return class
func (h NftClass) ClassByID(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.NftClasses{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Id:         h.Id(ctx),
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.Show(ctx, params)
}

func (h *NftClass) TransferNftClassByID(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftClassByIDRequest)
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

	// 校验参数 end
	authData := h.AuthData(ctx)
	params := dto.TransferNftClassById{
		ClassID:     h.ClassID(ctx),
		Owner:       h.Owner(ctx),
		Recipient:   recipient,
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	return h.svc.Transfer(ctx, params)
}

func (h NftClass) Id(ctx context.Context) string {
	idR := ctx.Value("id")
	if idR == nil {
		return ""
	}
	return idR.(string)
}

func (h NftClass) Name(ctx context.Context) string {
	nameR := ctx.Value("name")
	if nameR == nil {
		return ""
	}
	return nameR.(string)
}

func (h NftClass) Owner(ctx context.Context) string {
	ownerR := ctx.Value("owner")
	if ownerR == nil {
		return ""
	}
	return ownerR.(string)
}

func (h NftClass) TxHash(ctx context.Context) string {
	txHashR := ctx.Value("tx_hash")
	if txHashR == nil {
		return ""
	}
	return txHashR.(string)
}

func (h *NftClass) ClassID(ctx context.Context) string {
	classID := ctx.Value("class_id")
	if classID == nil {
		return ""
	}
	return classID.(string)
}
