package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type MTMsgsController struct {
	BaseController
	handler handlers.IMTMsgs
}

func NewMTMsgsController(bc BaseController, handler handlers.IMTMsgs) kit.IController {
	return MTMsgsController{bc, handler}
}

func (c MTMsgsController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/mt/mts/{mt_class_id}/{mt_id}/history",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.GetMTHistory, nil),
		},
	)
	return ends
}
