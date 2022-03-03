package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	types2 "github.com/irisnet/core-sdk-go/types"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type INft interface {
	CreateNft(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByBatch(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error)
	Nfts(ctx context.Context, _ interface{}) (interface{}, error)
	NftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	NftOperationHistoryByNftId(ctx context.Context, _ interface{}) (interface{}, error)
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
	tagBytes, _ := json.Marshal(req.Tag)

	tag := string(tagBytes)

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

	params := dto.CreateNftsP{
		ChainId:   h.ChainID(ctx),
		ClassId:   h.ClassId(ctx),
		Name:      name,
		Uri:       uri,
		UriHash:   uriHash,
		Data:      data,
		Amount:    req.Amount,
		Recipient: recipient,
		Tag:       tagBytes,
	}

	if params.Amount == 0 {
		params.Amount = 1
	}
	if params.Amount > 100 || params.Amount < 0 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrAmountInt)
	}
	return h.svc.CreateNfts(params)
}

// EditNftByNftId Edit a nft and return the edited result
func (h nft) EditNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {

	req := request.(*vo.EditNftByIndexRequest)

	name := strings.TrimSpace(req.Name)
	uri := strings.TrimSpace(req.Uri)
	data := strings.TrimSpace(req.Data)
	tag := strings.TrimSpace(req.Tag)

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

	//check end
	params := dto.EditNftByNftIdP{
		ChainId: h.ChainID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		NftId:   h.NftId(ctx),

		Name: name,
		Uri:  uri,
		Data: data,
		Tag:  tag,
	}

	return h.svc.EditNftByNftId(params)
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

		if v.NftId == "" {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftIdLen+" or "+types.ErrNftIdString)
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
	hash := make(map[string]bool)
	for i, Nft := range params.EditNfts {
		if hash[Nft.NftId] == false {
			hash[Nft.NftId] = true
		} else {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrRepeat)
		}
	}

	//check end
	return h.svc.EditNftByBatch(params)
}

// DeleteNftByNftId Delete a nft and return the edited result
func (h nft) DeleteNftByNftId(ctx context.Context, _ interface{}) (interface{}, error) {
	//check start

	//check end
	params := dto.DeleteNftByNftIdP{
		ChainId: h.ChainID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		NftId:   h.NftId(ctx),
	}

	return h.svc.DeleteNftByNftId(params)
}

// DeleteNftByBatch Delete multiple nfts and return the deleted results
func (h nft) DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {

	// check start
	nftIds, err := h.NftIds(ctx)
	if err != nil {
		return nil, err
	}

	// 2.judge whether the NFT is repeated
	hash := make(map[string]bool)
	for i, nftId := range nftIds {
		if nftId == "" {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftIdString)
		}
		if hash[nftId] == false {
			hash[nftId] = true
		} else {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrRepeat)
		}
	}

	if len(nftIds) > 50 {
		return nil, types.ErrLimit
	}
	//check end

	params := dto.DeleteNftByBatchP{
		ChainId: h.ChainID(ctx),
		ClassId: h.ClassId(ctx),
		Sender:  h.Owner(ctx),
		NftIds:  nftIds,
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

// NftByNftId return class details
func (h nft) NftByNftId(ctx context.Context, _ interface{}) (interface{}, error) {

	//check start

	//check end
	params := dto.NftByNftIdP{
		ChainId: h.ChainID(ctx),
		ClassId: h.ClassId(ctx),
		NftId:   h.NftId(ctx),
	}

	return h.svc.NftByNftId(params)

}

// NftOperationHistoryByNftId return class details
func (h nft) NftOperationHistoryByNftId(ctx context.Context, _ interface{}) (interface{}, error) {

	// 校验参数 start
	//check start

	params := dto.NftOperationHistoryByNftIdP{
		ClassID: h.ClassId(ctx),
		NftId:   h.NftId(ctx),
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
	return h.svc.NftOperationHistoryByNftId(params)
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
func (h nft) NftId(ctx context.Context) string {
	nftId := ctx.Value("nft_id")
	if nftId == nil {
		return ""
	}
	return nftId.(string)
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

func (h nft) NftIds(ctx context.Context) ([]string, error) {
	rec := ctx.Value("nft_ids")
	if rec == nil || rec == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndicesLen)
	}

	// "1,2,3,4,..." to {1,2,3,4,...}
	var ids []string
	strArr := strings.Split(rec.(string), ",")
	for _, s := range strArr {
		ids = append(ids, s)
	}

	return ids, nil
}
