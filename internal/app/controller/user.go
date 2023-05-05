package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

type UserController struct {
	base.BaseController
	handler handlers.IUser
}

func NewUserController(bc base.BaseController, handler handlers.IUser) kit.IController {
	return UserController{bc, handler}
}

func (c UserController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/users",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateUsers, &vo.CreateUserRequest{}),
		},
		kit.Endpoint{
			URI:     "/users",
			Method:  http.MethodPatch,
			Handler: c.MakeHandler(c.handler.UpdateUsers, &vo.UpdateUserRequest{}),
		},
		kit.Endpoint{
			URI:     "/users",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.ShowUsers, nil),
		},
	)
	return ends
}
