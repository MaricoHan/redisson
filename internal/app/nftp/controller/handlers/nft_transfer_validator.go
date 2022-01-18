package handlers

import "context"

type NftTransferValidator struct {
	next INftTransfer
}

// TransferNftClassByID transfer an nft class by id
func (h NftTransferValidator) TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error) {

	return h.next.TransferNftClassByID(ctx, nil)
}

// TransferNftByIndex transfer an nft class by index
func (h NftTransferValidator) TransferNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.TransferNftByIndex(ctx, nil)
}

// TransferNftByBatch return class list
func (h NftTransferValidator) TransferNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.TransferNftByBatch(ctx, nil)
}
