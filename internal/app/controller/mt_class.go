package controller

import (
	"net/http"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
)

type MTClassController struct {
	BaseController
	handler handlers.IMTClass
}

func NewMTClassController(bc BaseController, handler handlers.IMTClass) kit.IController {
	return MTClassController{bc, handler}
}

func (m MTClassController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/mt/classes",
			Method:  http.MethodGet,
			Handler: m.makeHandler(m.handler.List, nil),
		},
		kit.Endpoint{
			URI:     "/mt/classes",
			Method:  http.MethodPost,
			Handler: m.makeHandler(m.handler.CreateMTClass, &vo.CreateMTClassRequest{}),
		},
		kit.Endpoint{
			URI:     "/mt/classes/{mt_class_id}",
			Method:  http.MethodGet,
			Handler: m.makeHandler(m.handler.Show, nil),
		},
	)
	return ends
}
