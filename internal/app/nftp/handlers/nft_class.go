package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"
	
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

type INftClass interface {
	CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error)
	Classes(ctx context.Context, _ interface{}) (interface{}, error)
	ClassByID(ctx context.Context, _ interface{}) (interface{}, error)
}

type nftClass struct {
	base
	pageBasic
	svc map[string]service.NFTClassService
}

func NewNFTClass(svc ...*service.NFTClassBase) *nftClass {
	return newNFTClassModule(svc)
}

func newNFTClassModule(svc []*service.NFTClassBase) *nftClass {
	modules := make(map[string]service.NFTClassService, len(svc))
	for _, v := range svc {
		modules[v.Module] = v.Service
	}
	return &nftClass{
		svc: modules,
	}
}

// CreateNftClass Create one nft class
// return creation result
func (h nftClass) CreateNftClass(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateNftClassRequest)

	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)
	symbol := strings.TrimSpace(req.Symbol)
	uri := strings.TrimSpace(req.Uri)
	uriHash := strings.TrimSpace(req.UriHash)
	data := strings.TrimSpace(req.Data)
	owner := strings.TrimSpace(req.Owner)

	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrName)
	}

	if len([]rune(name)) < 3 || len([]rune(name)) > 64 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrNameLen)
	}

	if (symbol != "" && len([]rune(symbol)) < 3) || len([]rune(symbol)) > 64 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSymbolLen)
	}

	if len([]rune(description)) > 2048 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrDescriptionLen)
	}

	if err := h.base.UriCheck(uri); err != nil {
		return nil, err
	}

	if len([]rune(uriHash)) > 512 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrUriHashLen)
	}

	if len([]rune(data)) > 4096 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrDataLen)
	}

	if owner == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOwner)
	}

	if len([]rune(owner)) > 128 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOwnerLen)
	}

	authData := h.AuthData(ctx)
	params := dto.CreateNftClassP{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Data:        data,
		Owner:       owner,
		Tag:         tagBytes,
	}
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Create(params)
}

// Classes return class list
func (h nftClass) Classes(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.NftClassesP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Id:         h.Id(ctx),
		Name:       h.Name(ctx),
		Owner:      h.Owner(ctx),
		TxHash:     h.TxHash(ctx),
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
		startDateTime, err := time.Parse(timeLayout, fmt.Sprintf("%s 00:00:00", startDateR))
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrStartDate)
		}
		params.StartDate = &startDateTime
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		endDateTime, err := time.Parse(timeLayout, fmt.Sprintf("%s 23:59:59", endDateR))
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
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.List(params)
}

// ClassByID return class
func (h nftClass) ClassByID(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.NftClassesP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Id:         h.Id(ctx),
	}

	// 校验参数 end
	// 业务数据入库的地方
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Show(params)
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
