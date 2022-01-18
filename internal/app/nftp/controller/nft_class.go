package controller

import (
	"context"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type NftClassController struct {
	BaseController
}

func NewNftClassController(bc BaseController) kit.IController {
	return NftClassController{bc}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (c NftClassController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.Classes, nil),
		},
		kit.Endpoint{
			URI:     "/nft/classes",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.CreateNftClass, nil),
		},
		kit.Endpoint{
			URI:     "/nft/classes/{id}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.ClassByID, nil),
		},
	)
	return ends
}

// CreateNftClass Create one or more nft class
// return creation result
func (c NftClassController) CreateNftClass(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// Classes return class list
func (c NftClassController) Classes(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}

// ClassByID return class list
func (c NftClassController) ClassByID(ctx context.Context, _ interface{}) (interface{}, error) {
	panic("not yet implemented")
}
