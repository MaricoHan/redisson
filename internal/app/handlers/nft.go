package handlers

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INft interface {
	CreateNft(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	Nfts(ctx context.Context, _ interface{}) (interface{}, error)
	NftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	BatchTransfer(ctx context.Context, _ interface{}) (interface{}, error)
	BatchEdit(ctx context.Context, _ interface{}) (interface{}, error)
	BatchDelete(ctx context.Context, _ interface{}) (interface{}, error)
}

type NFT struct {
	base
	pageBasic
	svc services.INFT
}

func NewNft(svc services.INFT) *NFT {
	return &NFT{svc: svc}
}

func (h *NFT) CreateNft(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.CreateNftsRequest)

	name := strings.TrimSpace(req.Name)
	uri := strings.TrimSpace(req.Uri)
	uriHash := strings.TrimSpace(req.UriHash)
	data := strings.TrimSpace(req.Data)
	recipient := strings.TrimSpace(req.Recipient)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return nil, errors2.New(errors2.ClientParams, constant.ErrName)
	}

	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	if err := h.base.UriCheck(uri); err != nil {
		return nil, err
	}

	authData := h.AuthData(ctx)
	params := dto.CreateNfts{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		ClassId:    h.ClassId(ctx),
		Name:       name,
		Uri:        uri,
		UriHash:    uriHash,
		Data:       data,
		//Amount:    req.Amount,
		Recipient:   recipient,
		Tag:         tagBytes,
		Code:        authData.Code,
		OperationId: operationId,
	}

	params.Amount = 1

	return h.svc.Create(params)
}

// EditNftByNftId Edit a nft and return the edited result
func (h *NFT) EditNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {

	req := request.(*vo.EditNftByIndexRequest)

	name := strings.TrimSpace(req.Name)
	uri := strings.TrimSpace(req.Uri)
	data := strings.TrimSpace(req.Data)
	operationId := strings.TrimSpace(req.OperationID)
	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}

	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	//check start
	if name == "" {
		return nil, errors2.New(errors2.ClientParams, constant.ErrName)
	}

	if err := h.base.UriCheck(uri); err != nil {
		return nil, err
	}

	//check end
	authData := h.AuthData(ctx)
	params := dto.EditNftByNftId{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		ClassId:     h.ClassId(ctx),
		Sender:      h.Owner(ctx),
		NftId:       h.NftId(ctx),
		Module:      authData.Module,
		Name:        name,
		Uri:         uri,
		Data:        data,
		Tag:         tagBytes,
		Code:        authData.Code,
		OperationId: operationId,
	}

	return h.svc.Update(params)
}

// DeleteNftByNftId Delete a nft and return the edited result
func (h *NFT) DeleteNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {

	var tagBytes []byte
	var err error
	req := request.(*vo.DeleteNftByNftIdRequest)
	if req.Tag != nil {
		tagBytes, err = h.ValidateTag(req.Tag)
		if err != nil {
			return nil, err
		}
	}
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	authData := h.AuthData(ctx)
	params := dto.DeleteNftByNftId{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		ClassId:     h.ClassId(ctx),
		Sender:      h.Owner(ctx),
		NftId:       h.NftId(ctx),
		Tag:         tagBytes,
		Code:        authData.Code,
		OperationId: operationId,
	}

	return h.svc.Delete(params)
}

// Nfts return nft list
func (h *NFT) Nfts(ctx context.Context, _ interface{}) (interface{}, error) {
	status, err := h.Status(ctx)
	if err != nil {
		return nil, err
	}
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := dto.Nfts{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Id:         h.Id(ctx),
		ClassId:    h.ClassId(ctx),
		Owner:      h.Owner(ctx),
		TxHash:     h.TxHash(ctx),
		Status:     status,
		Code:       authData.Code,
		Name:       h.Name(ctx),
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
		params.StartDate = startDateR
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {

		params.EndDate = endDateR
	}

	params.SortBy = h.SortBy(ctx)
	// 校验参数 end
	// 业务数据入库的地方

	return h.svc.List(params)
}

// NftByNftId return class details
func (h *NFT) NftByNftId(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)
	params := dto.NftByNftId{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		ClassId:    h.ClassId(ctx),
		NftId:      h.NftId(ctx),
		Code:       authData.Code,
	}

	return h.svc.Show(params)

}

func (h *NFT) BatchTransfer(ctx context.Context, request interface{}) (interface{}, error) {
	// 接收请求
	req, ok := request.(*vo.BatchTransferRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 校验tag
	var tagBz []byte
	var err error
	if req.Tag != nil {
		tagBz, err = h.ValidateTag(req.Tag)
		if err != nil {
			return nil, err
		}
	}
	// 校验operationID
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	// 获取账户基本信息
	authData := h.AuthData(ctx)
	params := dto.BatchTransferRequest{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
		Sender:      h.Owner(ctx),
		Data:        req.Data,
		Tag:         string(tagBz),
		OperationID: operationId,
	}

	return h.svc.BatchTransfer(&params)
}

func (h *NFT) BatchEdit(ctx context.Context, request interface{}) (interface{}, error) {
	// 接收请求
	req, ok := request.(*vo.BatchEditRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 校验tag
	var tagBz []byte
	var err error
	if req.Tag != nil {
		tagBz, err = h.ValidateTag(req.Tag)
		if err != nil {
			return nil, err
		}
	}
	// 校验operationID
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	// 获取账户基本信息
	authData := h.AuthData(ctx)
	params := dto.BatchEditRequest{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
		Sender:      h.Owner(ctx),
		Nfts:        req.Nfts,
		Tag:         string(tagBz),
		OperationID: operationId,
	}

	return h.svc.BatchEdit(&params)
}

func (h *NFT) BatchDelete(ctx context.Context, request interface{}) (interface{}, error) {
	// 接收请求
	req, ok := request.(*vo.BatchDeleteRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 校验tag
	var tagBz []byte
	var err error
	if req.Tag != nil {
		tagBz, err = h.ValidateTag(req.Tag)
		if err != nil {
			return nil, err
		}
	}
	// 校验operationID
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len(operationId) == 0 || len(operationId) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	// 获取账户基本信息
	authData := h.AuthData(ctx)
	params := dto.BatchDeleteRequest{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
		Sender:      h.Owner(ctx),
		Nfts:        req.Nfts,
		Tag:         string(tagBz),
		OperationID: operationId,
	}

	return h.svc.BatchDelete(&params)
}

func (h *NFT) Signer(ctx context.Context) string {
	signer := ctx.Value("signer")
	if signer == nil || signer == "" {
		return ""
	}
	return signer.(string)
}

func (h *NFT) Operation(ctx context.Context) string {
	operation := ctx.Value("operation")
	if operation == nil || operation == "" {
		return ""
	}
	return operation.(string)
}

func (h *NFT) Txhash(ctx context.Context) string {
	txhash := ctx.Value("tx_hash")
	if txhash == nil || txhash == "" {
		return ""
	}
	return txhash.(string)
}

func (h *NFT) Id(ctx context.Context) string {
	id := ctx.Value("id")

	if id == nil {
		return ""
	}
	return id.(string)

}

func (h *NFT) ClassId(ctx context.Context) string {
	classId := ctx.Value("class_id")

	if classId == nil {
		return ""
	}
	return classId.(string)

}

func (h *NFT) Name(ctx context.Context) string {
	name := ctx.Value("name")
	if name == nil {
		return ""
	}
	return name.(string)
}

func (h *NFT) Owner(ctx context.Context) string {
	owner := ctx.Value("owner")
	if owner == nil {
		return ""
	}
	return owner.(string)

}
func (h *NFT) NftId(ctx context.Context) string {
	nftId := ctx.Value("nft_id")
	if nftId == nil {
		return ""
	}
	return nftId.(string)
}
func (h *NFT) TxHash(ctx context.Context) string {
	txHash := ctx.Value("tx_hash")
	if txHash == nil {
		return ""
	}

	return txHash.(string)
}
func (h *NFT) Status(ctx context.Context) (string, error) {
	status := ctx.Value("status")
	if status == nil || status == "" {
		return constant.NFTSStatusActive, nil
	}
	if status != constant.NFTSStatusActive && status != constant.NFTSStatusBurned {
		return "", errors2.New(errors2.ClientParams, errors2.ErrStatus)
	}
	return status.(string), nil
}
