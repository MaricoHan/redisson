package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type RightsController struct {
	BaseController
	handler handlers.IRights
}

func NewRightsController(bc BaseController, handler handlers.IRights) kit.IController {
	return RightsController{bc, handler}
}

func (r RightsController) GetEndpoints() []kit.Endpoint {
	return []kit.Endpoint{
		{
			URI:     "/rights/register",
			Method:  http.MethodPost,
			Handler: r.makeHandler(r.handler.Register, &vo.RegisterRequest{}),
		},
		{
			URI:     "/rights/register/{operation_id}",
			Method:  http.MethodPost,
			Handler: r.makeHandler(r.handler.EditRegister, &vo.EditRegisterRequest{}),
		},
		{
			URI:     "/rights/register/{register_type}/{operation_id}",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.QueryRegister, nil),
		},
	}
}
