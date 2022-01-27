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
func (h nft) CreateNft(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateNftsRequest)
	params := dto.CreateNftsRequest{
		AppID:     h.AppID(ctx),
		ClassId:   h.ClassId(ctx),
		Name:      req.Name,
		Uri:       req.Uri,
		UriHash:   req.UriHash,
		Data:      req.Data,
		Amount:    req.Amount,
		Recipient: req.Recipient,
	}
	if params.Amount == 0 {
		params.Amount = 1
	}
	if params.Amount > 100 {
		return nil, types.ErrParams
	}

	return h.svc.CreateNfts(params)
}

// EditNftByIndex Edit an nft and return the edited result
func (h nft) EditNftByIndex(ctx context.Context, request interface{}) (interface{}, error) {

	req := request.(*vo.EditNftByIndexRequest)

	//check start
	index, err := h.Index(ctx)
	if err != nil {
		return nil, err
	}
	//check end

	params := dto.EditNftByIndexP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		Index:   index,

		Name: req.Name,
		Uri:  req.Uri,
		Data: req.Data,
	}

	return h.svc.EditNftByIndex(params)
}

// EditNftByBatch Edit multiple nfts and
// return the deleted results
func (h nft) EditNftByBatch(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.EditNftByBatchRequest)
	params := dto.EditNftByBatchP{
		EditNfts: req.EditNftsR,
		AppID:    h.AppID(ctx),
		ClassId:  h.ClassId(ctx),
		Sender:   h.Owner(ctx),
	}

	//check start
	//1. count limit :50
	if len(params.EditNfts) > 50 {
		return nil, types.ErrNftTooMany
	}

	//check end

	return h.svc.EditNftByBatch(params)
}

// DeleteNftByIndex Delete an nft and return the edited result
func (h nft) DeleteNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {

	//check start
	index, err := h.Index(ctx)
	if err != nil {
		return nil, err
	}
	//check end
	params := dto.DeleteNftByIndexP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		Index:   index,
	}

	return h.svc.DeleteNftByIndex(params)
}

// DeleteNftByBatch Delete multiple nfts and
// return the deleted results
func (h nft) DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {

	// check start
	indices, err := h.Indices(ctx)
	if err != nil {
		return nil, types.ErrIndicesFormat
	}
	if len(indices) > 50 {
		return nil, types.ErrNftTooMany
	}
	//check end

	params := dto.DeleteNftByBatchP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		Indices: indices,
	}

	return h.svc.DeleteNftByBatch(params)
}

// Nfts return class list
func (h nft) Nfts(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	params := dto.NftsP{
		AppID:   h.AppID(ctx),
		Id:      h.Id(ctx),
		ClassId: h.ClassId(ctx),
		Owner:   h.Owner(ctx),
		TxHash:  h.TxHash(ctx),
		Status:  h.Status(ctx),
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

	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit >= 50 {
		return nil, types.ErrParams
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
	return h.svc.Nfts(params)
}

// NftByIndex return class details
func (h nft) NftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {

	//check start
	index, err := h.Index(ctx)
	if err != nil {
		return nil, err
	}
	//check end
	params := dto.NftByIndexP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Index:   index,
	}

	return h.svc.NftByIndex(params)

}

// NftOperationHistoryByIndex return class details
func (h nft) NftOperationHistoryByIndex(ctx context.Context, _ interface{}) (interface{}, error) {

	// 校验参数 start
	//check start
	index, err := h.Index(ctx)
	if err != nil {
		return nil, err
	}
	params := dto.NftOperationHistoryByIndexP{
		ClassID: h.ClassId(ctx),
		Index:   index,
		AppID:   h.AppID(ctx),
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

	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit >= 50 {
		return nil, types.ErrParams
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

	params.Signer = h.Signer(ctx)
	params.Txhash = h.Txhash(ctx)
	params.Operation = h.Operation(ctx)
	if params.Operation != "" {
		if !strings.Contains("mint/edit/transfer/burn", params.Operation) {
			return nil, types.ErrParams
		}
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

func (h nft) Id(ctx context.Context) string {
	id := ctx.Value("id")

	if id == nil {
		return ""
	}
	return id.(string)

}

func (h nft) ClassId(ctx context.Context) string {
	classId := ctx.Value("class_id")

	if classId == nil {
		return ""
	}
	return classId.(string)

}

func (h nft) Owner(ctx context.Context) string {
	owner := ctx.Value("owner")

	if owner == nil {
		return ""
	}
	return owner.(string)

}
func (h nft) Index(ctx context.Context) (uint64, error) {
	rec := ctx.Value("index")
	if rec == nil {
		return 0, types.ErrIndexFormat
	}
	index, err := strconv.ParseUint(rec.(string), 10, 64)
	if err == nil {
		return 0, types.ErrIndexFormat
	}

	// return index
	return index, nil
}
func (h nft) TxHash(ctx context.Context) string {
	txHash := ctx.Value("tx_hash")
	if txHash == nil {
		return ""
	}

	return txHash.(string)
}
func (h nft) Status(ctx context.Context) string {
	status := ctx.Value("status")
	if status == nil {
		return ""
	}
	return status.(string)
}

func (h nft) Indices(ctx context.Context) ([]uint64, error) {
	rec := ctx.Value("indices")
	if rec == nil {
		return nil, types.ErrIndicesFormat
	}

	// "1,2,3,4,..." to {1,2,3,4,...}
	var indices []uint64
	strArr := strings.Split(rec.(string), ",")

	for _, s := range strArr {
		tmp, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		indices = append(indices, tmp)
	}

	return indices, nil
}
