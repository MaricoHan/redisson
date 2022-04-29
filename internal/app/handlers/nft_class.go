package handlers

import (
	"context"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"strings"
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
		return nil, errors2.New(errors2.ClientParams, errors2.ErrName)
	}

	if len([]rune(name)) < 3 || len([]rune(name)) > 64 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrNameLen)
	}

	if (symbol != "" && len([]rune(symbol)) < 3) || len([]rune(symbol)) > 64 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSymbolLen)
	}

	if len([]rune(description)) > 2048 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrDescriptionLen)
	}

	if err := h.base.UriCheck(uri); err != nil {
		return nil, err
	}

	if len([]rune(uriHash)) > 512 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrURIHashLen)
	}

	if len([]rune(data)) > 4096 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrDataLen)
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
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Data:        data,
		Owner:       owner,
		Tag:         tagBytes,
		Code:        authData.Code,
	}
	return h.svc.CreateNFTClass(params)
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
	//switch h.SortBy(ctx) {
	//case "DATE_ASC":
	//	params.SortBy = "DATE_ASC"
	//case "DATE_DESC":
	//	params.SortBy = "DATE_DESC"
	//default:
	//	return nil, constant.NewAppError(constant.RootCodeSpace, constant.ClientParamsError, constant.ErrSortBy)
	//}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.GetAllNFTClasses(params)
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
	}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.GetNFTClass(params)
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
