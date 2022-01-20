package handlers

import (
	"context"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
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
	req := request.(vo.TransferNftClassByID)
	params := dto.TransferNftClassByIDP{
		Recipient: req.Recipient,
		AppID:     h.AppID(ctx),
	}

	// 校验参数 end
	return h.svc.TransferNftClassByID(params), nil
}

// TransferNftByIndex transfer an nft class by index
func (h nftTransfer) TransferNftByIndex(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(vo.TransferNftByIndex)

	params := dto.TransferNftByIndexP{
		Recipient: req.Recipient,
		AppID:     h.AppID(ctx),
	}
	// 校验参数 end
	return h.svc.TransferNftByIndex(params), nil
}

// TransferNftByBatch return class list
func (h nftTransfer) TransferNftByBatch(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(vo.TransferNftByBatch)
	params := dto.TransferNftByBatchP{
		Recipients: req.Recipients,
		AppID:      h.AppID(ctx),
	}
	// 校验参数 end
	return h.svc.TransferNftByBatch(params), nil
}
