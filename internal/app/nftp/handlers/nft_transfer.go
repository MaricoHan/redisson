package handlers

import (
	"context"
	"encoding/json"
	"strings"

	types2 "github.com/irisnet/core-sdk-go/types"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

type INftTransfer interface {
	TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftByBatch(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewNftTransfer(svc *service.NftTransfer) INftTransfer {
	return newNftTransfer(svc)
}

type nftTransfer struct {
	base
	svc *service.NftTransfer
}

func newNftTransfer(svc *service.NftTransfer) *nftTransfer {
	return &nftTransfer{svc: svc}
}

// TransferNftClassByID transfer an nft class by id
func (h nftTransfer) TransferNftClassByID(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftClassByIDRequest)
	recipient := strings.TrimSpace(req.Recipient)
	if recipient == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipient)
	}
	if len([]rune(recipient)) > 128 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientLen)
	}
	// 校验接收者地址是否满足当前链的地址规范
	if err := types2.ValidateAccAddress(recipient); err != nil {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
	}

	var tagBytes []byte
	if len(req.Tag) > 0 {
		tagBytes, _ := json.Marshal(req.Tag)
		tag := string(tagBytes)
		if _, err := h.IsValTag(tag); err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, err.Error())
		}
	}

	//校验参数 end

	params := dto.TransferNftClassByIDP{
		ClassID:    h.ClassID(ctx),
		Owner:      h.Owner(ctx),
		Recipient:  recipient,
		ChainID:    h.ChainID(ctx),
		ProjectID:  h.ProjectID(ctx),
		PlatFormID: h.PlatFormID(ctx),
		Tag:        tagBytes,
	}
	return h.svc.TransferNftClassByID(params)
}

// TransferNftByNftId transfer an nft class by index
func (h nftTransfer) TransferNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftByNftIdRequest)
	recipient := strings.TrimSpace(req.Recipient)
	if recipient == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipient)
	}
	if len([]rune(recipient)) > 128 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientLen)
	}
	// 校验接收者地址是否满足当前链的地址规范
	if err := types2.ValidateAccAddress(recipient); err != nil {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
	}
	var tagBytes []byte
	if len(req.Tag) > 0 {
		tagBytes, _ := json.Marshal(req.Tag)
		tag := string(tagBytes)
		if _, err := h.IsValTag(tag); err != nil {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, err.Error())
		}
	}

	// 校验参数 end

	params := dto.TransferNftByNftIdP{
		ClassID:    h.ClassID(ctx),
		Owner:      h.Owner(ctx),
		NftId:      h.NftId(ctx),
		Recipient:  recipient,
		ChainID:    h.ChainID(ctx),
		ProjectID:  h.ProjectID(ctx),
		PlatFormID: h.PlatFormID(ctx),
		Tag:        tagBytes,
	}
	return h.svc.TransferNftByNftId(params)
}

// TransferNftByBatch return class list
func (h nftTransfer) TransferNftByBatch(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftByBatchRequest)
	if req.Recipients == nil {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipients)
	}
	params := dto.TransferNftByBatchP{
		ClassID:    h.ClassID(ctx),
		Owner:      h.Owner(ctx),
		Recipients: req.Recipients,
		ChainID:    h.ChainID(ctx),
		ProjectID:  h.ProjectID(ctx),
		PlatFormID: h.PlatFormID(ctx),
	}
	if len(params.Recipients) > 50 {
		return "", types.ErrLimit
	}
	// 校验参数 end
	return h.svc.TransferNftByBatch(params)
}

func (h nftTransfer) ClassID(ctx context.Context) string {
	class_id := ctx.Value("class_id")
	if class_id == nil {
		return ""
	}
	return class_id.(string)
}

func (h nftTransfer) Owner(ctx context.Context) string {
	owner := ctx.Value("owner")
	if owner == nil {
		return ""
	}
	return owner.(string)
}

func (h nftTransfer) NftId(ctx context.Context) string {
	nftId := ctx.Value("nft_id")
	if nftId == nil {
		return ""
	}
	return nftId.(string)
}
