package controller

import (
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/handlers"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
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
			URI:     "/nft/classes",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.Classes, nil),
		},
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateNftClass, nil),
		},
		kit.Endpoint{
			URI:     "/nft/classes/{id}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.ClassByID, nil),
		},
	)
	return ends
}
