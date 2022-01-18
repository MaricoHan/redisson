package controller

import (
	"context"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type AccountController struct {
	BaseController
}

func NewAccountsController(bc BaseController) kit.IController {
	return AccountController{bc}
}

// GetEndpoints implement the method GetRouter of the interface IController
func (c AccountController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/accounts",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.CreateAccount, nil),
		},
		kit.Endpoint{
			URI:     "/accounts",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.Accounts, nil),
		},
	)
	return ends
}

// CreateAccount Create one or more accounts
// return creation result
func (c AccountController) CreateAccount(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// Accounts return account list
func (c AccountController) Accounts(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
