package native

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/nft"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/native"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INftClass interface {
	CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error)
	Classes(ctx context.Context, _ interface{}) (interface{}, error)
	ClassByID(ctx context.Context, _ interface{}) (interface{}, error)
}

type NftClass struct {
	handlers.Base
	handlers.PageBasic
	svc native.INFTClass
}

func NewNFTClass(svc native.INFTClass) *NftClass {
	return &NftClass{svc: svc}
}

// CreateNftClass Create one nft class
// return creation result
func (h NftClass) CreateNftClass(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateNftClassRequest)

	name := strings.TrimSpace(req.Name)
	classId := strings.TrimSpace(req.ClassId)
	description := strings.TrimSpace(req.Description)
	symbol := strings.TrimSpace(req.Symbol)
	uri := strings.TrimSpace(req.Uri)
	uriHash := strings.TrimSpace(req.UriHash)
	data := strings.TrimSpace(req.Data)
	owner := strings.TrimSpace(req.Owner)
	operationID := strings.TrimSpace(req.OperationID)

	if operationID == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if name == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrName)
	}
	if len(operationID) == 0 || len(operationID) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if err := h.UriCheck(uri); err != nil {
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
		ChainID:         authData.ChainId,
		ProjectID:       authData.ProjectId,
		PlatFormID:      authData.PlatformId,
		Module:          authData.Module,
		Name:            name,
		Symbol:          symbol,
		Description:     description,
		Uri:             uri,
		UriHash:         uriHash,
		Data:            data,
		Owner:           owner,
		Code:            authData.Code,
		OperationId:     operationID,
		ClassId:         classId,
		AccessMode:      authData.AccessMode,
		EditableByOwner: req.EditableByOwner,
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

	params.CountTotal, err = h.CountTotal(ctx)
	if err != nil {
		return nil, err
	}

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
