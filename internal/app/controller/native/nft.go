package native

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/nft"
)

type NftController struct {
	base.BaseController
	handler native.INft
}

func NewNftController(bc base.BaseController, handler native.INft) kit.IController {
	return NftController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateNft, &vo.CreateNftsRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/batch/nfts/{class_id}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.BatchCreateNft, &vo.BatchCreateNftsRequest{}),
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
		kit.Endpoint{
			URI:     "/nft/batch/nft-transfers/{owner}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.BatchTransfer, &vo.BatchTransferRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/batch/nfts/{owner}",
			Method:  http.MethodPatch,
			Handler: c.MakeHandler(c.handler.BatchEdit, &vo.BatchEditRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/batch/nfts/{owner}",
			Method:  http.MethodDelete,
			Handler: c.MakeHandler(c.handler.BatchDelete, &vo.BatchDeleteRequest{}),
		},
	)
	return ends
}
