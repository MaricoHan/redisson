package handlers

import (
	"context"
	"strconv"
	"strings"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

type INftTransfer interface {
	TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftByIndex(ctx context.Context, _ interface{}) (interface{}, error)
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
	params := dto.TransferNftClassByIDP{
		ClassID:   h.ClassID(ctx),
		Owner:     h.Owner(ctx),
		Recipient: recipient,
		AppID:     h.AppID(ctx),
	}
	//校验参数 end
	return h.svc.TransferNftClassByID(params)
}

// TransferNftByIndex transfer an nft class by index
func (h nftTransfer) TransferNftByIndex(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftByIndexRequest)
	recipient := strings.TrimSpace(req.Recipient)
	if recipient == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipient)
	}

	index, err := h.Index(ctx)
	if err != nil {
		return nil, err
	}

	if index == 0 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrIndexInt)
	}

	if len([]rune(recipient)) > 128 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientLen)
	}
	params := dto.TransferNftByIndexP{
		ClassID:   h.ClassID(ctx),
		Owner:     h.Owner(ctx),
		Index:     index,
		Recipient: recipient,
		AppID:     h.AppID(ctx),
	}
	// 校验参数 end
	return h.svc.TransferNftByIndex(params)
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
		AppID:      h.AppID(ctx),
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

func (h nftTransfer) Index(ctx context.Context) (uint64, error) {
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
