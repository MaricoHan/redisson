package controller

import (
	"net/http"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

type NFTTransferController struct {
	BaseController
	handler handlers.INFTTransfer
}

func NewNftTransferController(bc BaseController, handler handlers.INFTTransfer) kit.IController {
	return NFTTransferController{bc, handler}
}

func (c NFTTransferController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/class-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferNftClassByID, &vo.TransferNftClassByIDRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/nft-transfers/{class_id}/{owner}/{nft_id}",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferNftByNftId, &vo.TransferNftByNftIdRequest{}),
		},
	)
	return ends
}
