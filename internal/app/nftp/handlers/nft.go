package handlers

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"strings"
	"time"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type INft interface {
	CreateNft(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	Nfts(ctx context.Context, _ interface{}) (interface{}, error)
	NftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	NftOperationHistoryByNftId(ctx context.Context, _ interface{}) (interface{}, error)
}

type nft struct {
	base
	pageBasic
	svc map[string]service.NFTService
}

func NewNft(svc ...*service.NFTBase) INft {
	return newNFTModule(svc)
}

func newNFTModule(svc []*service.NFTBase) *nft {
	modules := make(map[string]service.NFTService, len(svc))
	for _, v := range svc {
		modules[v.Module] = v.Service
	}
	return &nft{
		svc: modules,
	}
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
		if !common.IsHexAddress(recipient){
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
		}
	}
	authData := h.AuthData(ctx)
	params := dto.CreateNftsP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		ClassId:    h.ClassId(ctx),
		Name:       name,
		Uri:        uri,
		UriHash:    uriHash,
		Data:       data,
		//Amount:    req.Amount,
		Recipient: recipient,
		Tag:       tagBytes,
	}
	if params.Amount == 0 {
		params.Amount = 1
	}
	if params.Amount > 100 || params.Amount < 0 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrAmountInt)
	}
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Create(params)
}

// EditNftByNftId Edit a nft and return the edited result
func (h nft) EditNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {

	req := request.(*vo.EditNftByIndexRequest)

	name := strings.TrimSpace(req.Name)
	uri := strings.TrimSpace(req.Uri)
	data := strings.TrimSpace(req.Data)
	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}
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
	authData := h.AuthData(ctx)
	params := dto.EditNftByNftIdP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		ClassId:    h.ClassId(ctx),
		Sender:     h.Owner(ctx),
		NftId:      h.NftId(ctx),

		Name: name,
		Uri:  uri,
		Data: data,
		Tag:  tagBytes,
	}

	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Update(params)
}

// DeleteNftByNftId Delete a nft and return the edited result
func (h nft) DeleteNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {

	var tagBytes []byte
	var err error
	if request != nil {
		req := request.(*vo.DeleteNftByNftIdRequest)
		tagBytes, err = h.ValidateTag(req.Tag)
		if err != nil {
			return nil, err
		}
	}
	authData := h.AuthData(ctx)
	params := dto.DeleteNftByNftIdP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		ClassId:    h.ClassId(ctx),
		Sender:     h.Owner(ctx),
		NftId:      h.NftId(ctx),
		Tag:        tagBytes,
	}

	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Delete(params)
}

// Nfts return class list
func (h nft) Nfts(ctx context.Context, _ interface{}) (interface{}, error) {
	status, err := h.Status(ctx)
	if err != nil {
		return nil, err
	}
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.NftsP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Id:         h.Id(ctx),
		ClassId:    h.ClassId(ctx),
		Owner:      h.Owner(ctx),
		TxHash:     h.TxHash(ctx),
		Status:     status,
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
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.List(params)
}

// NftByNftId return class details
func (h nft) NftByNftId(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)
	params := dto.NftByNftIdP{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		ClassId:    h.ClassId(ctx),
		NftId:      h.NftId(ctx),
	}
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.Show(params)

}

// NftOperationHistoryByNftId return class details
func (h nft) NftOperationHistoryByNftId(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)
	params := dto.NftOperationHistoryByNftIdP{
		ClassID:    h.ClassId(ctx),
		NftId:      h.NftId(ctx),
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
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
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.History(params)
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
