package handlers

import "context"

type INftTransfer interface {
	TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftByIndex(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftByBatch(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewNftTransfer() INftTransfer {
	return &NftTransferValidator{
		next: newNftTransfer(),
	}
}

type nftTransfer struct {
}

func newNftTransfer() *nftTransfer {
	return &nftTransfer{}
}

// TransferNftClassByID transfer an nft class by id
func (h nftTransfer) TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// TransferNftByIndex transfer an nft class by index
func (h nftTransfer) TransferNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// TransferNftByBatch return class list
func (h nftTransfer) TransferNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
