package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type MtClassController struct {
	BaseController
	handler handlers.IMtClass
}

func NewMtClassController(bc BaseController, handler handlers.IMtClass) kit.IController {
	return MtClassController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interfce IController
func (m MtClassController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/mt/classes",
			Method:  http.MethodGet,
			Handler: m.makeHandler(m.handler.List, nil),
		},
		kit.Endpoint{
			URI:     "/mt/classes/{mt_class_id}",
			Method:  http.MethodGet,
			Handler: m.makeHandler(m.handler.Show, nil),
		},
	)
	return ends
}
