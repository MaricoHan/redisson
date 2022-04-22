package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type AccountController struct {
	BaseController
	handler handlers.IAccount
}

func NewAccountsController(bc BaseController, handler handlers.IAccount) kit.IController {
	return AccountController{bc, handler}
}

func (c AccountController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/accounts",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateAccount, &vo.CreateAccountRequest{}),
		},
		kit.Endpoint{
			URI:     "/accounts",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.GetAccounts, nil),
		},

	)
	return ends
}
