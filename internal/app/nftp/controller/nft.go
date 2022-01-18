package controller

import (
	"context"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type NftController struct {
	BaseController
}

func NewNftController(bc BaseController) kit.IController {
	return NftController{bc}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.CreateNft, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}/{index}",
			Method:  http.MethodPatch,
			Handler: c.makeHandler(c.EditNftByIndex, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}",
			Method:  http.MethodPatch,
			Handler: c.makeHandler(c.EditNftByBatch, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}/{index}",
			Method:  http.MethodDelete,
			Handler: c.makeHandler(c.DeleteNftByIndex, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}",
			Method:  http.MethodPatch,
			Handler: c.makeHandler(c.DeleteNftByBatch, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.Nfts, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{index}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.NftByIndex, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{index}/history",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.NftOperationHistoryByIndex, nil),
		},
	)
	return ends
}

// CreateNft Create one or more nft class
// return creation result
func (c NftController) CreateNft(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// EditNftByIndex Edit an nft and return the edited result
func (c NftController) EditNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// EditNftByBatch Edit multiple nfts and
// return the deleted results
func (c NftController) EditNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// DeleteNftByIndex Delete an nft and return the edited result
func (c NftController) DeleteNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// DeleteNftByBatch Delete multiple nfts and
// return the deleted results
func (c NftController) DeleteNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// Nfts return class list
func (c NftController) Nfts(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// NftByIndex return class details
func (c NftController) NftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// NftOperationHistoryByIndex return class details
func (c NftController) NftOperationHistoryByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
