package native

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	native_handlers "gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
)

type TxController struct {
	base.BaseController
	handler native_handlers.ITx
}

func NewTxController(bc base.BaseController, handler native_handlers.ITx) kit.IController {
	return TxController{bc, handler}
}

func (c TxController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/native/tx/{operation_id}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.TxResult, nil),
		},
		//kit.Endpoint{
		//	URI:     "/tx/queue/info",
		//	Method:  http.MethodGet,
		//	Handler: c.MakeHandler(c.handler.TxQueueInfo, nil),
		//},
	)
	return ends
}
