package controller

import (
	"context"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
)

type DemoController struct {
	BaseController
}

func NewDemoController(bc BaseController) kit.IController {
	return DemoController{bc}
}

// GetEndpoints implement the method GetRouter of the interface IController
func (d DemoController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends, kit.Endpoint{
		URI:     "/demo",
		Method:  http.MethodGet,
		Handler: d.makeHandler(d.Demo, nil),
	})
	return ends
}

// Demo return a demo
func (d DemoController) Demo(ctx context.Context, _ interface{}) (interface{}, error) {
	ctx.Value("X-App-ID")
	return map[string]string{"demo": "this is demo"}, nil
}
