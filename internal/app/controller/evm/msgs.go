package evm

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/evm"
)

type MsgsController struct {
	base.BaseController
	handler evm.IMsgs
}

func NewMsgsController(bc base.BaseController, handler evm.IMsgs) kit.IController {
	return MsgsController{bc, handler}
}

func (c MsgsController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/evm/nft/nfts/{class_id}/{nft_id}/history",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.GetNFTHistory, nil),
		},
	)
	return ends
}
