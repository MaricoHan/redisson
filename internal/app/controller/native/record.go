package native

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	vo "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/record"
	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
)

type RecordController struct {
	base.BaseController
	handler native.IRecord
}

func NewRecordController(bc base.BaseController, handler native.IRecord) kit.IController {
	return RecordController{bc, handler}
}

// GetEndpoints implement the method GetRouter of the interface IController
func (c RecordController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/native/record/records",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.CreateRecord, &vo.RecordCreateRequest{}),
		},
	)
	return ends
}
