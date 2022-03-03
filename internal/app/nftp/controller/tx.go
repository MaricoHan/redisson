package controller

import (
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/handlers"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type TxController struct {
	BaseController
	handler handlers.ITx
}

func NewTxController(bc BaseController, handler handlers.ITx) kit.IController {
	return TxController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
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
