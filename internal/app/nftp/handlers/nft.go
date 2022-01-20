package handlers

import (
	"context"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
)

type INft interface {
	CreateNft(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByIndex(ctx context.Context, _ interface{}) (interface{}, error)
	EditNftByBatch(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByIndex(ctx context.Context, _ interface{}) (interface{}, error)
	DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error)
	Nfts(ctx context.Context, _ interface{}) (interface{}, error)
	NftByIndex(ctx context.Context, _ interface{}) (interface{}, error)
	NftOperationHistoryByIndex(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewNft(svc *service.Nft) INft {
	return newNft(svc)
}

type nft struct {
	base
	svc *service.Nft
}

func newNft(svc *service.Nft) *nft {
	return &nft{svc: svc}
}

// CreateNft Create one or more nft class
// return creation result
func (h nft) CreateNft(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")

}

// EditNftByIndex Edit an nft and return the edited result
func (h nft) EditNftByIndex(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(vo.EditNftByIndexRequest)
	params := dto.EditNftByIndexP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Owner:   h.Owner(ctx),
		Index:   h.Index(ctx),

		Name: req.Name,
		Uri:  req.Uri,
		Data: req.Data,
	}
	//check start
	//1. judge whether the Caller is the owner

	//2. judge whether the Caller is the APP's address

	//check end
	return h.svc.EditNftByIndex(params)
}

// EditNftByBatch Edit multiple nfts and
// return the deleted results
func (h nft) EditNftByBatch(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(vo.EditNftByBatchRequest)

	params := dto.EditNftByBatchP{

		Index: req.Index,
		Name:  req.Name,
		Uri:   req.Uri,
		Data:  req.Data,
	}

	return h.svc.EditNftByBatch(params)
}

// DeleteNftByIndex Delete an nft and return the edited result
func (h nft) DeleteNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	h.Index(ctx)
	panic("not yet implemented")
}

// DeleteNftByBatch Delete multiple nfts and
// return the deleted results
func (h nft) DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// Nfts return class list
func (h nft) Nfts(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// NftByIndex return class details
func (h nft) NftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	params := dto.NftByIndexP{
		AppID:   h.AppID(ctx),
		ClassId: h.ClassId(ctx),
		Index:   h.Index(ctx),
	}
	//check start

	//check end
	return h.svc.NftByIndex(params)

}

// NftOperationHistoryByIndex return class details
func (h nft) NftOperationHistoryByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

func (h nft) ClassId(ctx context.Context) string {
	class_id := ctx.Value("class_id")

	if class_id == nil {
		return ""
	}
	return class_id.(string)

}

func (h nft) Owner(ctx context.Context) string {
	owner := ctx.Value("owner")

	if owner == nil {
		return ""
	}
	return owner.(string)

}
func (h nft) Index(ctx context.Context) uint64 {
	index := ctx.Value("index")

	return index.(uint64)
}
