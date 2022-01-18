package handlers

import "context"

type nftValidator struct {
	next INft
}

// CreateNft Create one or more nft class
// return creation result
func (h nftValidator) CreateNft(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.CreateNft(ctx, nil)
}

// EditNftByIndex Edit an nft and return the edited result
func (h nftValidator) EditNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.EditNftByIndex(ctx, nil)
}

// EditNftByBatch Edit multiple nfts and
// return the deleted results
func (h nftValidator) EditNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.EditNftByBatch(ctx, nil)
}

// DeleteNftByIndex Delete an nft and return the edited result
func (h nftValidator) DeleteNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.DeleteNftByIndex(ctx, nil)
}

// DeleteNftByBatch Delete multiple nfts and
// return the deleted results
func (h nftValidator) DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.DeleteNftByBatch(ctx, nil)
}

// Nfts return class list
func (h nftValidator) Nfts(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.Nfts(ctx, nil)
}

// NftByIndex return class details
func (h nftValidator) NftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.NftByIndex(ctx, nil)
}

// NftOperationHistoryByIndex return class details
func (h nftValidator) NftOperationHistoryByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.NftOperationHistoryByIndex(ctx, nil)
}
