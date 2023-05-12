package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
)

type AuthController struct {
	base.BaseController
	handler handlers.IAuth
}

func NewAuthController(bc base.BaseController, handler handlers.IAuth) kit.IController {
	return AuthController{bc, handler}
}

func (c AuthController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/auth/verify",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.Verify, nil),
		},
		kit.Endpoint{
			URI:     "/auth/users",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.GetUser, nil),
		},
	)
	return ends
}
