package controller

import (
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
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
			Handler: c.makeHandler(c.handler.TransferNftClassByID, &vo.TransferNftClassByIDRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nft-transfers/{class_id}/{owner}/{index}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferNftByIndex, &vo.TransferNftByIndexRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nft-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferNftByBatch, &vo.TransferNftByBatchRequest{}),
		},
	)
	return ends
}
