package handlers

import (
	"context"
	"fmt"
	types2 "github.com/irisnet/core-sdk-go/types"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
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

	name := strings.TrimSpace(req.Name)
	uri := strings.TrimSpace(req.Uri)
	uriHash := strings.TrimSpace(req.UriHash)
	data := strings.TrimSpace(req.Data)
	recipient := strings.TrimSpace(req.Recipient)
	tag := strings.TrimSpace(req.Tag)

	if name == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrName)
	}

	if len([]rune(name)) < 3 || len([]rune(name)) > 64 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrNameLen)
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

	if len([]rune(recipient)) > 128 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientLen)
	}
	// 若接收者地址不为空，则校验其格式；否则在service中将其默认设为NFT类别的权属者地址
	if recipient != "" {
		// 校验接收者地址是否满足当前链的地址规范
		if err := types2.ValidateAccAddress(recipient); err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
		}
	}

	params := dto.CreateNftsRequest{
		ChainId:   h.ChainID(ctx),
		ClassId:   h.ClassId(ctx),
		Name:      name,
		Uri:       uri,
		UriHash:   uriHash,
		Data:      data,
		Amount:    req.Amount,
		Recipient: recipient,
		Tag:       tag,
	}

	if params.Amount == 0 {
		params.Amount = 1
	}
	if params.Amount > 100 || params.Amount < 0 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrAmountInt)
	}
	return h.svc.CreateNfts(params)
}

// EditNftByIndex Edit an nft and return the edited result
func (h nft) EditNftByIndex(ctx context.Context, request interface{}) (interface{}, error) {

	req := request.(*vo.EditNftByIndexRequest)

	name := strings.TrimSpace(req.Name)
	uri := strings.TrimSpace(req.Uri)
	data := strings.TrimSpace(req.Data)

	//check start
	if name == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrName)
	}
	if len([]rune(name)) < 3 || len([]rune(name)) > 64 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrNameLen)
	}
	if err := h.base.UriCheck(uri); err != nil {
		return nil, err
	}
	if len([]rune(data)) > 4096 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrDataLen)
	}

	index, err := h.Index(ctx)
	if err != nil {
		return nil, err
	}
	if index == 0 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndexInt)
	}
	//check end

	params := dto.EditNftByIndexP{
		ChainId: h.ChainID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		Index:   index,

		Name: name,
		Uri:  uri,
		Data: data,
	}

	return h.svc.EditNftByIndex(params)
}

// EditNftByBatch Edit multiple nfts and
// return the deleted results
func (h nft) EditNftByBatch(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.EditNftByBatchRequest)

	if len(*req) == 0 {
		return "", nil
	}

	var nfts []*dto.EditNft
	for i, v := range *req {
		v.Name = strings.TrimSpace(v.Name)
		v.Uri = strings.TrimSpace(v.Uri)
		v.Data = strings.TrimSpace(v.Data)

		if v.Index == 0 {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrIndexLen+" or "+types.ErrIndexInt)
		}

		if v.Name == "" {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrName)
		}

		if len([]rune(v.Name)) < 3 || len([]rune(v.Name)) > 64 {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNameLen)
		}

		if v.Uri != "" {

			if len([]rune(v.Uri)) > 256 {
				return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrUriLen)
			}

			isUri := govalidator.IsRequestURI(v.Uri)
			if !isUri {
				return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrUri)
			}
		}

		if len([]rune(v.Data)) > 4096 {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrDataLen)
		}

		nfts = append(nfts, v)
	}

	params := dto.EditNftByBatchP{
		EditNfts: nfts,
		ChainId:  h.ChainID(ctx),
		ClassId:  h.ClassId(ctx),
		Sender:   h.Owner(ctx),
	}
	//check start
	//1. count limit :50
	if len(params.EditNfts) > 50 {
		return nil, types.ErrLimit
	}
	// 2.judge whether the NFT is repeated
	hash := make(map[uint64]bool)
	for i, Nft := range params.EditNfts {
		if hash[Nft.Index] == false {
			hash[Nft.Index] = true
		} else {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrRepeat)
		}
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
	if index == 0 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndexInt)
	}

	//check end
	params := dto.DeleteNftByIndexP{
		ChainId: h.ChainID(ctx),
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
		return nil, err
	}

	// 2.judge whether the NFT is repeated
	hash := make(map[uint64]bool)
	for i, index := range indices {
		if index == 0 {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrIndexInt)
		}
		if hash[index] == false {
			hash[index] = true
		} else {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrRepeat)
		}
	}

	if len(indices) > 50 {
		return nil, types.ErrLimit
	}
	//check end

	params := dto.DeleteNftByBatchP{
		ChainId: h.ChainID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		Indices: indices,
	}

	return h.svc.DeleteNftByBatch(params)
}

// Nfts return class list
func (h nft) Nfts(ctx context.Context, _ interface{}) (interface{}, error) {
	status, err := h.Status(ctx)
	if err != nil {
		return nil, err
	}
	// 校验参数 start
	params := dto.NftsP{
		ChainId: h.ChainID(ctx),
		Id:      h.Id(ctx),
		ClassId: h.ClassId(ctx),
		Owner:   h.Owner(ctx),
		TxHash:  h.TxHash(ctx),
		Status:  status,
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
	case "ID_ASC":
		params.SortBy = "ID_ASC"
	case "ID_DESC":
		params.SortBy = "ID_DESC"
	case "DATE_ASC":
		params.SortBy = "DATE_ASC"
	case "DATE_DESC":
		params.SortBy = "DATE_DESC"
	default:
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSortBy)
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
	if index == 0 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndexInt)
	}

	//check end
	params := dto.NftByIndexP{
		ChainId: h.ChainID(ctx),
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
	if index == 0 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndexInt)
	}

	params := dto.NftOperationHistoryByIndexP{
		ClassID: h.ClassId(ctx),
		Index:   index,
		ChainId: h.ChainID(ctx),
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

	params.Signer = h.Signer(ctx)
	params.Txhash = h.Txhash(ctx)
	params.Operation = h.Operation(ctx)
	if params.Operation != "" {
		if !strings.Contains("mint/edit/transfer/burn", params.Operation) {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOperation)
		}
	}
	// 校验参数 end
	return h.svc.NftOperationHistoryByIndex(params)
}

func (h nft) Signer(ctx context.Context) string {
	signer := ctx.Value("signer")
	if signer == nil || signer == "" {
		return ""
	}
	return signer.(string)
}

func (h nft) Operation(ctx context.Context) string {
	operation := ctx.Value("operation")
	if operation == nil || operation == "" {
		return ""
	}
	return operation.(string)
}

func (h nft) Txhash(ctx context.Context) string {
	txhash := ctx.Value("tx_hash")
	if txhash == nil || txhash == "" {
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
		return 0, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndexLen)
	}
	index, err := strconv.ParseUint(rec.(string), 10, 64)
	if err != nil {
		return 0, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndex)
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
func (h nft) Status(ctx context.Context) (string, error) {
	status := ctx.Value("status")
	if status == nil || status == "" {
		return models.TNFTSStatusActive, nil
	}
	if status != models.TNFTSStatusActive && status != models.TNFTSStatusBurned {
		return "", types.ErrNftStatus
	}
	return status.(string), nil
}

func (h nft) Indices(ctx context.Context) ([]uint64, error) {
	rec := ctx.Value("indices")
	if rec == nil || rec == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndicesLen)
	}

	// "1,2,3,4,..." to {1,2,3,4,...}
	var indices []uint64
	strArr := strings.Split(rec.(string), ",")
	for _, s := range strArr {
		tmp, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndices)
		}
		indices = append(indices, tmp)
	}

	return indices, nil
}
