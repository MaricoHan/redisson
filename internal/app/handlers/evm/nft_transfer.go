package evm

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/evm"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/evm"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/evm"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"
)

type INFTTransfer interface {
	TransferNftClassByID(ctx context.Context, request interface{}) (interface{}, error)
	TransferNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
}

type NFTTransfer struct {
	base.Base
	base.PageBasic
	NFT
	svc evm.INFTTransfer
}

func NewNFTTransfer(svc evm.INFTTransfer) *NFTTransfer {
	return &NFTTransfer{svc: svc}
}

func (h *NFTTransfer) TransferNftClassByID(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftClassByIDRequest)
	recipient := strings.TrimSpace(req.Recipient)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if recipient == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrRecipient)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	// 校验参数 end
	authData := h.AuthData(ctx)
	params := dto.TransferNftClassById{
		ClassID:     h.ClassID(ctx),
		Owner:       h.Owner(ctx),
		Recipient:   recipient,
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	return h.svc.TransferNFTClass(ctx, params)
}

// TransferNftByNftId transfer a nft class by index
func (h *NFTTransfer) TransferNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftByNftIdRequest)
	recipient := strings.TrimSpace(req.Recipient)
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if recipient == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrRecipient)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	// 校验参数 end
	authData := h.AuthData(ctx)
	params := dto.TransferNftByNftId{
		ClassID:     h.ClassID(ctx),
		Sender:      h.Owner(ctx),
		Recipient:   recipient,
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}

	nftId, err := h.NftID(ctx)
	if err != nil {
		return nil, err
	}

	params.NftId = nftId

	// 不能自己转让给自己
	// 400
	if params.Recipient == params.Sender {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSelfTransfer)
	}

	return h.svc.TransferNFT(ctx, params)
}
