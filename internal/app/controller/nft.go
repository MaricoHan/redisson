package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
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
			Handler: c.makeHandler(c.handler.CreateNft, &vo.CreateNftsRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/batch/nfts/{class_id}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.BatchCreateNft, &vo.BatchCreateNftsRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodPatch,
			Handler: c.makeHandler(c.handler.EditNftByNftId, &vo.EditNftByIndexRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodDelete,
			Handler: c.makeHandler(c.handler.DeleteNftByNftId, &vo.DeleteNftByNftIdRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nfts",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.Nfts, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{nft_id}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.NftByNftId, nil),
		},
	)
	return ends
}
