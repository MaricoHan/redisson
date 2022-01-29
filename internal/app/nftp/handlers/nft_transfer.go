package handlers

import (
	"context"
	"strconv"

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
	if req.Recipient == "" {
		return nil, types.ErrParams
	}
	if req.Recipient != "" && len(req.Recipient) > 128 {
		return nil, types.ErrParams
	}
	params := dto.TransferNftClassByIDP{
		ClassID:   h.ClassID(ctx),
		Owner:     h.Owner(ctx),
		Recipient: req.Recipient,
		AppID:     h.AppID(ctx),
	}
	//校验参数 end
	return h.svc.TransferNftClassByID(params)
}

// TransferNftByIndex transfer an nft class by index
func (h nftTransfer) TransferNftByIndex(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftByIndexRequest)
	if req.Recipient == "" || h.Index(ctx) == 0 {
		return nil, types.ErrParams
	}

	if req.Recipient != "" && len(req.Recipient) > 128 {
		return nil, types.ErrParams
	}
	params := dto.TransferNftByIndexP{
		ClassID:   h.ClassID(ctx),
		Owner:     h.Owner(ctx),
		Index:     h.Index(ctx),
		Recipient: req.Recipient,
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
		return nil, types.ErrParams
	}
	params := dto.TransferNftByBatchP{
		ClassID:    h.ClassID(ctx),
		Owner:      h.Owner(ctx),
		Recipients: req.Recipients,
		AppID:      h.AppID(ctx),
	}
	if len(params.Recipients) > 50 {
		return "", types.ErrParams
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

func (h nftTransfer) Index(ctx context.Context) uint64 {
	index := ctx.Value("index")
	if index == nil {
		return 0
	}
	parseUint, err := strconv.ParseUint(index.(string), 10, 64)
	if err != nil {
		panic(err)
	}
	return parseUint
}
