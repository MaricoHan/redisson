package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

type NFTTransferController struct {
	base.BaseController
	handler handlers.INFTTransfer
}

func NewNftTransferController(bc base.BaseController, handler handlers.INFTTransfer) kit.IController {
	return NFTTransferController{bc, handler}
}

func (c NFTTransferController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/evm/nft/class-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferNftClassByID, &vo.TransferNftClassByIDRequest{}),
		},
		kit.Endpoint{
			URI:     "/evm/nft/nft-transfers/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferNftByNftId, &vo.TransferNftByNftIdRequest{}),
		},
		// 兼容之前的
		kit.Endpoint{
			URI:     "/nft/class-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferNftClassByID, &vo.TransferNftClassByIDRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nft-transfers/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferNftByNftId, &vo.TransferNftByNftIdRequest{}),
		},
	)
	return ends
}
