package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type MsgsController struct {
	BaseController
	handler handlers.IMsgs
}

func NewMsgsController(bc BaseController, handler handlers.IMsgs) kit.IController {
	return MsgsController{bc, handler}
}

func (c MsgsController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{nft_id}/history",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.GetNFTHistory, nil),
		},
		kit.Endpoint{
			URI:     "/accounts/history",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.GetAccountHistory, nil),
		},
	)
	return ends
}
