package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
)

type TxController struct {
	base.BaseController
	handler handlers.ITx
}

func NewTxController(bc base.BaseController, handler handlers.ITx) kit.IController {
	return TxController{bc, handler}
}

func (c TxController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/tx/{operation_id}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.TxResult, nil),
		},
		//kit.Endpoint{
		//	URI:     "/tx/queue/info",
		//	Method:  http.MethodGet,
		//	Handler: c.MakeHandler(c.handler.TxQueueInfo, nil),
		//},
	)
	return ends
}
