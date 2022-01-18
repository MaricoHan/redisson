package controller

import (
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/controller/handlers"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type DemoController struct {
	bc      BaseController
	handler handlers.IDemo
}

func NewDemoController(bc BaseController, handler handlers.IDemo) kit.IController {
	return DemoController{bc: bc, handler: handler}
}

// GetEndpoints implement the method GetRouter of the interface IController
func (d DemoController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/demo",
			Method:  http.MethodGet,
			Handler: d.bc.makeHandler(d.handler.Demo, nil),
		},
		kit.Endpoint{
			URI:     "/demo/{id}",
			Method:  http.MethodGet,
			Handler: d.bc.makeHandler(d.handler.DemoByID, nil),
		},
	)
	return ends
}
