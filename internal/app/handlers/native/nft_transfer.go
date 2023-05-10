package native

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/nft"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/native"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INFTTransfer interface {
	TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
}

type NFTTransfer struct {
	handlers.Base
	handlers.PageBasic
	svc native.INFTTransfer
}

func NewNFTTransfer(svc native.INFTTransfer) *NFTTransfer {
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

// TransferNftByNftId transfer an nft class by index
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
		NftId:       h.NftId(ctx),
		Recipient:   recipient,
		ChainID:     authData.ChainId,
		ProjectID:   authData.ProjectId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
		OperationId: operationId,
		AccessMode:  authData.AccessMode,
	}
	// 不能自己转让给自己
	// 400
	if params.Recipient == params.Sender {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSelfTransfer)
	}

	return h.svc.TransferNFT(ctx, params)
}

func (h *NFTTransfer) ClassID(ctx context.Context) string {
	classId := ctx.Value("class_id")
	if classId == nil {
		return ""
	}
	return classId.(string)
}

func (h *NFTTransfer) Owner(ctx context.Context) string {
	owner := ctx.Value("owner")
	if owner == nil {
		return ""
	}
	return owner.(string)
}

func (h *NFTTransfer) NftId(ctx context.Context) string {
	nftId := ctx.Value("nft_id")
	if nftId == nil {
		return ""
	}
	return nftId.(string)
}
