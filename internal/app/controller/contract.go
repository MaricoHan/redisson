package controller

import (
	"net/http"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

type ContractController struct {
	BaseController
	handler handlers.IContract
}

func NewContractController(bc BaseController, handler handlers.IContract) kit.IController {
	return ContractController{bc, handler}
}

func (c ContractController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/evm/contract/calls",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.ShowCall, nil),
		},

		kit.Endpoint{
			URI:     "/evm/contract/calls",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateCall, &vo.CreateContractCallRequest{}),
		},
		// 兼容之前的
		kit.Endpoint{
			URI:     "/contract/calls",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.ShowCall, nil),
		},

		kit.Endpoint{
			URI:     "/contract/calls",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateCall, &vo.CreateContractCallRequest{}),
		},
	)
	return ends
}
