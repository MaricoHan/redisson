package controller

import (
	"context"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type NftTransferController struct {
	BaseController
}

func NewNftTransferController(bc BaseController) kit.IController {
	return NftTransferController{bc}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftTransferController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/class-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.TransferNftClassByID, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nft-transfers/{class_id}/{owner}/{index}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.TransferNftByIndex, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nft-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.TransferNftByBatch, nil),
		},
	)
	return ends
}

// TransferNftClassByID transfer an nft class by id
func (c NftTransferController) TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// TransferNftByIndex transfer an nft class by index
func (c NftTransferController) TransferNftByIndex(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// TransferNftByBatch return class list
func (c NftTransferController) TransferNftByBatch(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
