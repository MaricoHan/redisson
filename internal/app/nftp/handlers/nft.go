package handlers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

type INft interface {
	CreateNft(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByIndex(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByBatch(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByIndex(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error)
	Nfts(ctx context.Context, _ interface{}) (interface{}, error)
	NftByIndex(ctx context.Context, _ interface{}) (interface{}, error)
	NftOperationHistoryByIndex(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewNft(svc *service.Nft) INft {
	return newNft(svc)
}

type nft struct {
	base
	pageBasic
	svc *service.Nft
}

func newNft(svc *service.Nft) *nft {
	return &nft{svc: svc}
}

// CreateNft Create one or more nft class
// return creation result
func (h nft) CreateNft(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")

}

// EditNftByIndex Edit an nft and return the edited result
func (h nft) EditNftByIndex(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(vo.EditNftByIndexRequest)
	params := dto.EditNftByIndexP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		Index:   h.Index(ctx),

		Name: req.Name,
		Uri:  req.Uri,
		Data: req.Data,
	}
	//check start
	//1. judge whether the Caller is the owner

	//2. judge whether the Caller is the APP's address

	//check end
	return h.svc.EditNftByIndex(params)
}

// EditNftByBatch Edit multiple nfts and
// return the deleted results
func (h nft) EditNftByBatch(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(vo.EditNftByBatchRequest)

	params := dto.EditNftByBatchP{
		EditNfts: req.EditNftsR,
		AppID:    h.AppID(ctx),
		ClassId:  h.ClassId(ctx),
		Sender:   h.Owner(ctx),
	}
	//check start

	//1. count limit :50
	if len(params.EditNfts) > 50 {
		return nil, types.ErrNftParams
	}

	//2. judge whether the Caller is the owner

	//3. judge whether the Caller is the APP's address

	//check end
	return h.svc.EditNftByBatch(params)
}

// DeleteNftByIndex Delete an nft and return the edited result
func (h nft) DeleteNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	params := dto.DeleteNftByIndexP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		Index:   h.Index(ctx),
	}
	//check start
	//1. judge whether the Caller is the owner

	//2. judge whether the Caller is the APP's address

	//check end

	return h.svc.DeleteNftByIndex(params)
}

// DeleteNftByBatch Delete multiple nfts and
// return the deleted results
func (h nft) DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {

	params := dto.DeleteNftByBatchP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		Indices: h.Indices(ctx),
	}

	//check start
	//1. judge whether the Caller is the owner

	//2. judge whether the Caller is the APP's address

	//check end

	return h.svc.DeleteNftByBatch(params)
}

// Nfts return class list
func (h nft) Nfts(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// NftByIndex return class details
func (h nft) NftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	params := dto.NftByIndexP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Index:   h.Index(ctx),
	}

	//check start
	//...
	//check end

	return h.svc.NftByIndex(params)

}

// NftOperationHistoryByIndex return class details
func (h nft) NftOperationHistoryByIndex(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	params := dto.NftOperationHistoryByIndexP{
		ClassID: h.ClassId(ctx),
		Index:   h.Index(ctx),
		AppID:   h.AppID(ctx),
	}
	//params.Signer = h.Signer(ctx)
	//params.Operation = h.Operation(ctx)
	//params.Txhash = h.Txhash(ctx)
	//
	//offset, err := h.Offset(ctx)
	//if err != nil {
	//	return nil, types.ErrParams
	//}
	//params.Offset = offset

	//limit, err := h.Limit(ctx)
	//if err != nil {
	//	return nil, types.ErrParams
	//}
	//params.Limit = limit

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
	return h.svc.NftOperationHistoryByIndex(params)
}

func (h nft) Signer(ctx context.Context) string {
	signer := ctx.Value("signer")
	if signer == nil {
		return ""
	}
	return signer.(string)
}

func (h nft) Operation(ctx context.Context) string {
	operation := ctx.Value("operation")
	if operation == nil {
		return ""
	}
	return operation.(string)
}

func (h nft) Txhash(ctx context.Context) string {
	txhash := ctx.Value("tx_hash")
	if txhash == nil {
		return ""
	}
	return txhash.(string)
}

func (h nft) ClassId(ctx context.Context) string {
	class_id := ctx.Value("class_id")

	if class_id == nil {
		return ""
	}
	return class_id.(string)

}

func (h nft) Owner(ctx context.Context) string {
	owner := ctx.Value("owner")

	if owner == nil {
		return ""
	}
	return owner.(string)

}
func (h nft) Index(ctx context.Context) uint64 {
	index := ctx.Value("index")

	if index == nil {
		return 0
	}
	parseUint, err := strconv.ParseUint(index.(string), 10, 64)
	if err != nil {
		panic(err)
	}
	return parseUint
}
func (h nft) Indices(ctx context.Context) []uint64 {
	rec := ctx.Value("indices")

	//"1,2,3,4,..." to {1,2,3,4,...}
	var indices []uint64
	strArr := strings.Split(rec.(string), ",")
	for i, s := range strArr {
		indices[i], _ = strconv.ParseUint(s, 10, 64)
	}

	return indices
}
