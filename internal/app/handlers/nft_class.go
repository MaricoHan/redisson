package handlers

import (
	"context"
	"fmt"
	"gitlab.bianjie.ai/avata/utils/errors/common"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INftClass interface {
	CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error)
	Classes(ctx context.Context, _ interface{}) (interface{}, error)
	ClassByID(ctx context.Context, _ interface{}) (interface{}, error)
}

type NftClass struct {
	base
	pageBasic
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
	symbol := strings.TrimSpace(req.Symbol)
	uri := strings.TrimSpace(req.Uri)
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
	if err := h.base.UriCheck(uri); err != nil {
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
		Module:      authData.Module,
		Name:        name,
		Symbol:      symbol,
		Uri:         uri,
		Owner:       owner,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	return h.svc.CreateNFTClass(ctx, params)
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
	return h.svc.GetAllNFTClasses(ctx, params)
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
	return h.svc.GetNFTClass(ctx, params)
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
