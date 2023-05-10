package native

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/mt"
)

type MTClassController struct {
	base.BaseController
	handler native.IMTClass
}

func NewMTClassController(bc base.BaseController, handler native.IMTClass) kit.IController {
	return MTClassController{bc, handler}
}

func (m MTClassController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/mt/classes",
			Method:  http.MethodGet,
			Handler: m.MakeHandler(m.handler.List, nil),
		},
		kit.Endpoint{
			URI:     "/mt/classes",
			Method:  http.MethodPost,
			Handler: m.MakeHandler(m.handler.CreateMTClass, &vo.CreateMTClassRequest{}),
		},
		kit.Endpoint{
			URI:     "/mt/classes/{id}",
			Method:  http.MethodGet,
			Handler: m.MakeHandler(m.handler.Show, nil),
		},
		kit.Endpoint{
			URI:     "/mt/class-transfers/{id}/{owner}",
			Method:  http.MethodPost,
			Handler: m.MakeHandler(m.handler.TransferMTClass, &vo.TransferMTClassRequest{}),
		},
	)
	return ends
}
