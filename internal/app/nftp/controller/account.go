package controller

import (
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/handlers"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type AccountController struct {
	BaseController
	handler handlers.IAccount
}

func NewAccountsController(bc BaseController, handler handlers.IAccount) kit.IController {
	return AccountController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interface IController
func (c AccountController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/accounts",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.Create, &vo.CreateAccountRequest{}),
		},
		kit.Endpoint{
			URI:     "/accounts",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.Accounts, nil),
		},
		kit.Endpoint{
			URI:     "/accounts/history",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.AccountsHistory, nil),
		},
	)
	return ends
}
