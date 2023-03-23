package controller

import (
	"net/http"

	"gitlab.bianjie.ai/avata/chains/api/pb/v1beta1/record"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
)

type RecordController struct {
	BaseController
	handler handlers.IRecord
}

func NewRecordController(bc BaseController, handler handlers.IRecord) kit.IController {
	return RecordController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interface IController
func (c RecordController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/record/records",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.CreateRecord, &record.RecordCreateRequest{}),
		},
	)
	return ends
}
