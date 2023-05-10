package native

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
)

type MsgsController struct {
	base.BaseController
	handler native.IMsgs
}

func NewMsgsController(bc base.BaseController, handler native.IMsgs) kit.IController {
	return MsgsController{bc, handler}
}

func (c MsgsController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/nfts/{class_id}/{nft_id}/history",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.GetNFTHistory, nil),
		},
		kit.Endpoint{
			URI:     "/accounts/history",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.GetAccountHistory, nil),
		},
		kit.Endpoint{
			URI:     "/mt/mts/{class_id}/{mt_id}/history",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.GetMTHistory, nil),
		},
	)
	return ends
}
