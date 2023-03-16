package controller

import (
	"net/http"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

type NsController struct {
	BaseController
	handler handlers.INs
}

func NewNsController(bc BaseController, handler handlers.INs) kit.IController {
	return NsController{bc, handler}
}

func (c NsController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/ns/domains",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.Domains, nil),
		},
		kit.Endpoint{
			URI:     "/ns/domains",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateDomain, &vo.CreateDomainRequest{}),
		},
	)
	return ends
}
