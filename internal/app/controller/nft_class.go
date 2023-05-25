package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/evm"
)

type NftClassController struct {
	base.BaseController
	handler evm.INftClass
}

func NewNftClassController(bc base.BaseController, handler evm.INftClass) kit.IController {
	return NftClassController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftClassController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		// 兼容之前的
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.Classes, nil),
		},
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateNftClass, &map[string]interface{}{}),
		},
		kit.Endpoint{
			URI:     "/nft/classes/{id}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.ClassByID, nil),
		},
	)
	return ends
}
