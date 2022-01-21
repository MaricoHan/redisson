package handlers

import (
	"context"
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
	if true {
		//参数校验
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
		return nil, types.ErrParams
	}
	params.Offset = offset

	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, types.ErrParams
	}
	params.Limit = limit
	if params.Offset == 0 {
		params.Offset = 1
	}

	if params.Limit == 0 {
		params.Limit = 10
	}
	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		startDateTime, err := time.Parse(timeLayout, startDateR)
		if err != nil {
			return nil, types.ErrParams
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, endDateR)
		if err != nil {
			return nil, types.ErrParams
		}
		params.EndDate = &endDateTime
	}

	if params.EndDate != nil && params.StartDate != nil {
		if !params.EndDate.After(*params.StartDate) {
			return nil, types.ErrParams
		}
	}
	switch h.SortBy(ctx) {
	case "DATE_ASC":
		params.SortBy = "DATE_ASC"
	case "DATE_DESC":
		params.SortBy = "DATE_DESC"
	default:
		return nil, types.ErrParams
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
