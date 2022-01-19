package handlers

import "context"

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

func NewNft() INft {
	return newNft()
}

type nft struct {
	base
}

func newNft() *nft {
	return &nft{}
}

// CreateNft Create one or more nft class
// return creation result
func (h nft) CreateNft(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// EditNftByIndex Edit an nft and return the edited result
func (h nft) EditNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// EditNftByBatch Edit multiple nfts and
// return the deleted results
func (h nft) EditNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// DeleteNftByIndex Delete an nft and return the edited result
func (h nft) DeleteNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
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
	panic("not yet implemented")
}

// NftOperationHistoryByIndex return class details
func (h nft) NftOperationHistoryByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
