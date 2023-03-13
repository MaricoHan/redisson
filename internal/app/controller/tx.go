package controller

import (
	"net/http"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

type TxController struct {
	BaseController
	handler handlers.ITx
}

func NewTxController(bc BaseController, handler handlers.ITx) kit.IController {
	return TxController{bc, handler}
}

func (c TxController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/tx/{operation_id}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.TxResult, nil),
		},
		kit.Endpoint{
			URI:     "/tx/queue/info",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.TxQueueInfo, nil),
		},
	)
	return ends
}
