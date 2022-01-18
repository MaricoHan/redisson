package handlers

import "context"

type txValidator struct {
	next ITx
}

// TxResultByTxHash transfer an nft class by id
func (h txValidator) TxResultByTxHash(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.TxResultByTxHash(ctx, nil)
}
