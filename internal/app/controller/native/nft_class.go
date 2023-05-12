package native

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/nft"
)

type NftClassController struct {
	base.BaseController
	handler native.INftClass
}

func NewNftClassController(bc base.BaseController, handler native.INftClass) kit.IController {
	return NftClassController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftClassController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.Classes, nil),
		},
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateNftClass, &vo.CreateNftClassRequest{}),
		},
		kit.Endpoint{
			URI:     "/nft/classes/{id}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.ClassByID, nil),
		},
	)
	return ends
}
