package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

type AccountController struct {
	base.BaseController
	handler handlers.IAccount
}

func NewAccountsController(bc base.BaseController, handler handlers.IAccount) kit.IController {
	return AccountController{bc, handler}
}

func (c AccountController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/accounts",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.BatchCreateAccount, &vo.BatchCreateAccountRequest{}),
		},
		kit.Endpoint{
			URI:     "/account",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateAccount, &vo.CreateAccountRequest{}),
		},
		kit.Endpoint{
			URI:     "/accounts",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.GetAccounts, nil),
		},
	)
	return ends
}
