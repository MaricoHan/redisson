package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

type NsController struct {
	base.BaseController
	handler handlers.INs
}

func NewNsController(bc base.BaseController, handler handlers.INs) kit.IController {
	return NsController{bc, handler}
}

func (c NsController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/evm/ns/domains",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.Domains, nil),
		},
		kit.Endpoint{
			URI:     "/evm/ns/domains/{owner}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.UserDomains, nil),
		},
		kit.Endpoint{
			URI:     "/evm/ns/domains",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateDomain, &vo.CreateDomainRequest{}),
		},
		kit.Endpoint{
			URI:     "/evm/ns/transfers/{owner}/{name}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferDomain, &vo.TransferDomainRequest{}),
		},
		// 兼容之前的
		kit.Endpoint{
			URI:     "/ns/domains",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.Domains, nil),
		},
		kit.Endpoint{
			URI:     "/ns/domains/{owner}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.UserDomains, nil),
		},
		kit.Endpoint{
			URI:     "/ns/domains",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateDomain, &vo.CreateDomainRequest{}),
		},
		kit.Endpoint{
			URI:     "/ns/transfers/{owner}/{name}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferDomain, &vo.TransferDomainRequest{}),
		},
	)
	return ends
}
