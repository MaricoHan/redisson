package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/nft_v2"

	"gitlab.bianjie.ai/avata/utils/errors/common"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INft interface {
	CreateNft(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	Nfts(ctx context.Context, _ interface{}) (interface{}, error)
	NftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
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

	uri := strings.TrimSpace(req.Uri)
	uriHash := strings.TrimSpace(req.UriHash)
	recipient := strings.TrimSpace(req.Recipient)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	if err := h.base.UriCheck(uri); err != nil {
		return nil, err
	}

	authData := h.AuthData(ctx)
	params := dto.CreateNfts{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		ClassId:     h.ClassId(ctx),
		Uri:         uri,
		UriHash:     uriHash,
		Recipient:   recipient,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.Create(ctx, params)
}

// EditNftByNftId Edit a nft and return the edited result
func (h *NFT) EditNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.EditNftByIndexRequest)

	uri := strings.TrimSpace(req.Uri)
	uriHash := strings.TrimSpace(req.UriHash)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	// check start

	if err := h.base.UriCheck(uri); err != nil {
		return nil, err
	}

	// check end
	authData := h.AuthData(ctx)
	params := dto.EditNftByNftId{
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		ClassId:     h.ClassId(ctx),
		Sender:      h.Owner(ctx),
		Module:      authData.Module,
		Uri:         uri,
		UriHash:     uriHash,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	nftId, err := h.NftId(ctx)
	if err != nil {
		return nil, err
	}

	params.NftId = nftId

	return h.svc.Update(ctx, params)
}

// DeleteNftByNftId Delete a nft and return the edited result
func (h *NFT) DeleteNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.DeleteNftByNftIdRequest)

	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
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
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	nftId, err := h.NftId(ctx)
	if err != nil {
		return nil, err
	}

	params.NftId = nftId

	return h.svc.Delete(ctx, params)
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
		Code:       authData.Code,
		AccessMode: authData.AccessMode,

		ClassId: h.ClassId(ctx),
		Owner:   h.Owner(ctx),
		TxHash:  h.TxHash(ctx),
		Status:  status,
	}
	params.Id, err = h.Id(ctx)
	if err != nil {
		return nil, err
	}

	params.PageKey = h.PageKey(ctx)
	countTotal, err := h.CountTotal(ctx)
	if err != nil {
		return nil, errors2.New(errors2.ClientParams, fmt.Sprintf(common.ERR_INVALID_VALUE, "count_total"))
	}
	params.CountTotal = countTotal

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

	return h.svc.List(ctx, params)
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
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}
	nftId, err := h.NftId(ctx)
	if err != nil {
		return nil, err
	}

	params.NftId = nftId

	return h.svc.Show(ctx, params)
}

func (h *NFT) Signer(ctx context.Context) string {
	signer := ctx.Value("signer")
	if signer == nil || signer == "" {
		return ""
	}
	return signer.(string)
}

func (h *NFT) Id(ctx context.Context) (uint64, error) {
	s := ctx.Value("id")
	if s == nil {
		return 0, nil
	}
	res, err := strconv.ParseUint(s.(string), 10, 64)
	if err != nil {
		return 0, errors2.New(errors2.ClientParams, fmt.Sprintf(common.ERR_INVALID_VALUE, "id"))
	}
	return res, nil
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

func (h *NFT) NftId(ctx context.Context) (uint64, error) {
	v := ctx.Value("nft_id")
	if v == nil {
		return 0, errors2.New(errors2.NotFound, "")
	}
	res, err := strconv.ParseUint(v.(string), 10, 64)
	if err != nil {
		return 0, errors2.New(errors2.NotFound, fmt.Sprintf("%s, nft_id: %s not found", errors2.ErrResourceNotFound, v.(string)))
	}

	return res, nil
}

func (h *NFT) TxHash(ctx context.Context) string {
	txHash := ctx.Value("tx_hash")
	if txHash == nil {
		return ""
	}

	return txHash.(string)
}

func (h *NFT) Status(ctx context.Context) (pb.STATUS, error) {
	v := ctx.Value("status")
	if v == nil {
		return pb.STATUS_ACTIVE, nil
	}
	s := v.(string)
	status, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, errors2.New(errors2.ClientParams, errors2.ErrStatus)
	}

	return pb.STATUS(status), nil
}
