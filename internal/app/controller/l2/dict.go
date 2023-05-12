package l2

import (
	"net/http"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/l2"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

type DictController struct {
	base.BaseController

	handler l2.IDict
}

var _ kit.IController = DictController{}

func NewDictController(b base.BaseController, h l2.IDict) kit.IController {
	return DictController{b, h}
}
func (d DictController) GetEndpoints() []kit.Endpoint {
	return []kit.Endpoint{
		{
			URI:     "/l2/dict/tx_types",
			Method:  http.MethodGet,
			Handler: d.MakeHandler(d.handler.ListTxTypes, nil),
		},
	}

}
