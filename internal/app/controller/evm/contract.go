package evm

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/evm"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/evm"
)

type ContractController struct {
	base.BaseController
	handler evm.IContract
}

func NewContractController(bc base.BaseController, handler evm.IContract) kit.IController {
	return ContractController{bc, handler}
}

func (c ContractController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/evm/contract/calls",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.ShowCall, nil),
		},

		kit.Endpoint{
			URI:     "/evm/contract/calls",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateCall, &vo.CreateContractCallRequest{}),
		},
	)
	return ends
}
