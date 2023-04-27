package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type NftClassController struct {
	BaseController
	handler handlers.INftClass
}

func NewNftClassController(bc BaseController, handler handlers.INftClass) kit.IController {
	return NftClassController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftClassController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/evm/nft/classes",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.Classes, nil),
		},
		kit.Endpoint{
			URI:     "/evm/nft/classes",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateNftClass, &vo.CreateNftClassRequest{}),
		},
		kit.Endpoint{
			URI:     "/evm/nft/classes/{id}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.ClassByID, nil),
		},
		// 兼容之前的
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.Classes, nil),
		},
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateNftClass, &vo.CreateNftClassRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/classes/{id}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.ClassByID, nil),
		},
	)
	return ends
}
