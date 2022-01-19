package handlers

import (
	"context"
)

type INftClass interface {
	CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error)
	Classes(ctx context.Context, _ interface{}) (interface{}, error)
	ClassByID(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewNftClass() INftClass {
	return newNftClass()
}

type nftClass struct {
	base
}

func newNftClass() *nftClass {
	return &nftClass{}
}

// CreateNftClass Create one or more nft class
// return creation result
func (h nftClass) CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// Classes return class list
func (h nftClass) Classes(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// ClassByID return class list
func (h nftClass) ClassByID(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
