package controller

import (
	"net/http"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

type UserController struct {
	BaseController
	handler handlers.IUser
}

func NewUserController(bc BaseController, handler handlers.IUser) kit.IController {
	return UserController{bc, handler}
}

func (c UserController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/users",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateUsers, &vo.CreateUserRequest{}),
		},
		kit.Endpoint{
			URI:     "/users",
			Method:  http.MethodPatch,
			Handler: c.makeHandler(c.handler.UpdateUsers, &vo.UpdateUserRequest{}),
		},
		kit.Endpoint{
			URI:     "/users",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.ShowUsers, &vo.ShowUserRequest{}),
		},
	)
	return ends
}
