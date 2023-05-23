package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/evm"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/evm"
)

type NftController struct {
	base.BaseController
	handler evm.INft
}

func NewNftController(bc base.BaseController, handler evm.INft) kit.IController {
	return NftController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		// 兼容之前的版本
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateNft, &vo.CreateNftsRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodPatch,
			Handler: c.MakeHandler(c.handler.EditNftByNftId, &vo.EditNftByIndexRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodDelete,
			Handler: c.MakeHandler(c.handler.DeleteNftByNftId, &vo.DeleteNftByNftIdRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nfts",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.Nfts, nil),
		},
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{nft_id}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.NftByNftId, nil),
		},
	)
	return ends
}
