package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

type AuthController struct {
	BaseController
	handler handlers.IAuth
}

func NewAuthController(bc BaseController, handler handlers.IAuth) kit.IController {
	return AuthController{bc, handler}
}

func (c AuthController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/auth/verify",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.Verify, &vo.AuthVerify{}),
		},
		kit.Endpoint{
			URI:     "/auth/users",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.GetUser, &vo.AuthGetUser{}),
		},
	)
	return ends
}
