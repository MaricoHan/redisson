package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
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
			URI:     "/tx/{task_id}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.TxResultByTxHash, nil),
		},
	)
	return ends
}
