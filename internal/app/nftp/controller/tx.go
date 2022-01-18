package controller

import (
	"context"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type TxController struct {
	BaseController
}

func NewTxController(bc BaseController) kit.IController {
	return TxController{bc}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c TxController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/tx/{hash}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.TxResultByTxHash, nil),
		},
	)
	return ends
}

// TxResultByTxHash transfer an nft class by id
func (c TxController) TxResultByTxHash(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
