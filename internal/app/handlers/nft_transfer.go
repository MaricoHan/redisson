package handlers

import (
	"context"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"strings"
)

type INFTTransfer interface {
	TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
}

type NFTTransfer struct {
	base
	pageBasic
	svc services.INFTTransfer
}

func NewNFTTransfer(svc services.INFTTransfer) *NFTTransfer {
	return &NFTTransfer{svc: svc}
}

func (h *NFTTransfer) TransferNftClassByID(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftClassByIDRequest)
	recipient := strings.TrimSpace(req.Recipient)
	if recipient == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrRecipient)
	}

	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}

	//校验参数 end
	authData := h.AuthData(ctx)
	params := dto.TransferNftClassById{
		ClassID:    h.ClassID(ctx),
		Owner:      h.Owner(ctx),
		Recipient:  recipient,
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Tag:        tagBytes,
		Code:       authData.Code,
	}
	return h.svc.TransferNFTClass(params)
}

// TransferNftByNftId transfer an nft class by index
func (h *NFTTransfer) TransferNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftByNftIdRequest)
	recipient := strings.TrimSpace(req.Recipient)
	if recipient == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrRecipient)
	}

	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}

	// 校验参数 end
	authData := h.AuthData(ctx)
	params := dto.TransferNftByNftId{
		ClassID:    h.ClassID(ctx),
		Sender:     h.Owner(ctx),
		NftId:      h.NftId(ctx),
		Recipient:  recipient,
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Module:     authData.Module,
		Tag:        tagBytes,
		Code:       authData.Code,
	}
	//不能自己转让给自己
	//400
	if params.Recipient == params.Sender {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSelfTransfer)
	}

	return h.svc.TransferNFT(params)
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
