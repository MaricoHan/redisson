package l2

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	handlers "gitlab.bianjie.ai/avata/open-api/internal/app/handlers/l2"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/l2"
)

type NftController struct {
	base.BaseController
	handler handlers.INft
}

func NewNftController(bc base.BaseController, handler handlers.INft) kit.IController {
	return NftController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/l2/nft/nfts/{class_id}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateNft, &vo.CreateNftsRequest{}),
		},
		kit.Endpoint{
			URI:     "/l2/nft/nfts/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodPatch,
			Handler: c.MakeHandler(c.handler.EditNftByNftId, &vo.EditNftByIndexRequest{}),
		},
		kit.Endpoint{
			URI:     "/l2/nft/nfts/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodDelete,
			Handler: c.MakeHandler(c.handler.DeleteNftByNftId, &vo.DeleteNftByNftIdRequest{}),
		},
		kit.Endpoint{
			URI:     "/l2/nft/nfts",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.Nfts, nil),
		},
		kit.Endpoint{
			URI:     "/l2/nft/nfts/{class_id}/{nft_id}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.NftByNftId, nil),
		},
		kit.Endpoint{
			URI:     "/l2/nft/nft-transfers/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferNftByNftId, &vo.TransferNftByNftIdRequest{}),
		},
	)
	return ends
}
