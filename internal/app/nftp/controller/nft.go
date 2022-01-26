package controller

import (
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/handlers"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type NftController struct {
	BaseController
	handler handlers.INft
}

func NewNftController(bc BaseController, handler handlers.INft) kit.IController {
	return NftController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateNft, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}/{index}",
			Method:  http.MethodPatch,
			Handler: c.makeHandler(c.handler.EditNftByIndex, &vo.EditNftByIndexRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}",
			Method:  http.MethodPatch,
			Handler: c.makeHandler(c.handler.EditNftByBatch, &vo.EditNftByBatchRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}/{index}",
			Method:  http.MethodDelete,
			Handler: c.makeHandler(c.handler.DeleteNftByIndex, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}",
			Method:  http.MethodDelete,
			Handler: c.makeHandler(c.handler.DeleteNftByBatch, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.Nfts, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{index}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.NftByIndex, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{index}/history",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.NftOperationHistoryByIndex, nil),
		},
	)
	return ends
}
