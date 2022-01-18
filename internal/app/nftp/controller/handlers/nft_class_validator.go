package handlers

import (
	"context"
)

type NftClassValidator struct {
	next INftClass
}

// CreateNftClass Create one or more nft class
// return creation result
func (h NftClassValidator) CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.CreateNftClass(ctx, nil)
}

// Classes return class list
func (h NftClassValidator) Classes(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.Classes(ctx, nil)
}

// ClassByID return class list
func (h NftClassValidator) ClassByID(ctx context.Context, _ interface{}) (interface{}, error) {
	return h.next.ClassByID(ctx, nil)
}
