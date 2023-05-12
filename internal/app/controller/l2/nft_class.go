package l2

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	handlers "gitlab.bianjie.ai/avata/open-api/internal/app/handlers/l2"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/l2"
)

type NftClassController struct {
	base.BaseController
	handler handlers.INftClass
}

func NewNftClassController(bc base.BaseController, handler handlers.INftClass) kit.IController {
	return NftClassController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftClassController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/l2/nft/classes",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.Classes, nil),
		},
		kit.Endpoint{
			URI:     "/l2/nft/classes",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateNftClass, &vo.CreateNftClassRequest{}),
		},
		kit.Endpoint{
			URI:     "/l2/nft/classes/{id}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.ClassByID, nil),
		},
		kit.Endpoint{
			URI:     "/l2/nft/class-transfers/{class_id}/{owner}",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferNftClassByID, &vo.TransferNftClassByIDRequest{}),
		},
	)
	return ends
}
