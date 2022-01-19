package handlers

import "context"

type ITx interface {
	TxResultByTxHash(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewTx() ITx {
	return newTx()
}

type tx struct {
	base
}

func newTx() *tx {
	return &tx{}
}

// TxResultByTxHash transfer an nft class by id
func (h tx) TxResultByTxHash(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
