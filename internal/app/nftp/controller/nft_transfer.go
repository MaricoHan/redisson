package controller

import (
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/handlers"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type NftTransferController struct {
	BaseController
	handler handlers.INftTransfer
}

func NewNftTransferController(bc BaseController, handler handlers.INftTransfer) kit.IController {
	return NftTransferController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftTransferController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/class-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferNftClassByID, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nft-transfers/{class_id}/{owner}/{index}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferNftByIndex, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nft-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferNftByBatch, nil),
		},
	)
	return ends
}
