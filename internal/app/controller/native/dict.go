package native

import (
	"net/http"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

type DictController struct {
	base.BaseController

	handler native.IDict
}

var _ kit.IController = DictController{}

func NewDictController(b base.BaseController, h native.IDict) kit.IController {
	return DictController{b, h}
}
func (d DictController) GetEndpoints() []kit.Endpoint {
	return []kit.Endpoint{
		{
			URI:     "/native/dict/tx_types",
			Method:  http.MethodGet,
			Handler: d.MakeHandler(d.handler.ListTxTypes, nil),
		},
	}

}
