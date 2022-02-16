package handlers

import (
	"context"
	"strings"
	"time"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

type INftClass interface {
	CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error)
	Classes(ctx context.Context, _ interface{}) (interface{}, error)
	ClassByID(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewNftClass(svc *service.NftClass) INftClass {
	return newNftClass(svc)
}

type nftClass struct {
	base
	pageBasic
	svc *service.NftClass
}

func newNftClass(svc *service.NftClass) *nftClass {
	return &nftClass{svc: svc}
}

// CreateNftClass Create one nft class
// return creation result
func (h nftClass) CreateNftClass(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateNftClassRequest)

	if req.Name == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrName)
	}

	if len([]rune(strings.TrimSpace(req.Name))) < 3 || len([]rune(strings.TrimSpace(req.Name))) > 64 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrNameLen)
	}

	if (req.Symbol != "" && len([]rune(strings.TrimSpace(req.Symbol))) != 0 && len([]rune(strings.TrimSpace(req.Symbol))) < 3) || (req.Symbol != "" && len([]rune(strings.TrimSpace(req.Symbol))) != 0 && len([]rune(strings.TrimSpace(req.Symbol))) > 64) {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSymbolLen)
	}

	if req.Description != "" && len([]rune(strings.TrimSpace(req.Description))) > 2048 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrDescriptionLen)
	}

	if err := h.base.UriCheck(&req.Uri); err != nil {
		return nil, err
	}

	if req.UriHash != "" && len([]rune(strings.TrimSpace(req.UriHash))) > 512 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrUriHashLen)
	}

	if req.Data != "" && len([]rune(strings.TrimSpace(req.Data))) > 4096 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrDataLen)
	}

	if req.Owner == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOwner)
	}

	if len([]rune(strings.TrimSpace(req.Owner))) > 128 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOwnerLen)
	}

	params := dto.CreateNftClassP{
		AppID:       h.AppID(ctx),
		Name:        req.Name,
		Symbol:      req.Symbol,
		Description: req.Description,
		Uri:         req.Uri,
		UriHash:     req.UriHash,
		Data:        req.Data,
		Owner:       req.Owner,
	}
	return h.svc.CreateNftClass(params)
}

// Classes return class list
func (h nftClass) Classes(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	params := dto.NftClassesP{
		AppID:  h.AppID(ctx),
		Id:     h.Id(ctx),
		Name:   h.Name(ctx),
		Owner:  h.Owner(ctx),
		TxHash: h.TxHash(ctx),
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
		startDateTime, err := time.Parse(timeLayoutWithoutHMS, startDateR)
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrStartDate)
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayoutWithoutHMS, endDateR)
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
	return h.svc.NftClasses(params)
}

// ClassByID return class
func (h nftClass) ClassByID(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	params := dto.NftClassesP{
		AppID: h.AppID(ctx),
		Id:    h.Id(ctx),
	}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.NftClassById(params)
}

func (h nftClass) Id(ctx context.Context) string {
	idR := ctx.Value("id")
	if idR == nil {
		return ""
	}
	return idR.(string)
}
func (h nftClass) Name(ctx context.Context) string {
	nameR := ctx.Value("name")
	if nameR == nil {
		return ""
	}
	return nameR.(string)
}
func (h nftClass) Owner(ctx context.Context) string {
	ownerR := ctx.Value("owner")
	if ownerR == nil {
		return ""
	}
	return ownerR.(string)
}
func (h nftClass) TxHash(ctx context.Context) string {
	txHashR := ctx.Value("tx_hash")
	if txHashR == nil {
		return ""
	}
	return txHashR.(string)
}
