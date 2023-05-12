package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

type EmptionController struct {
	base.BaseController
	handler handlers.IBusiness
}

func NewEmptionController(bc base.BaseController, handler handlers.IBusiness) kit.IController {
	return EmptionController{bc, handler}
}

func (c EmptionController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/orders/{operation_id}",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.GetOrderInfo, nil),
		},
		kit.Endpoint{
			URI:     "/orders",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.BuildOrder, &vo.BuyRequest{}),
		},
		kit.Endpoint{
			URI:     "/orders",
			Method:  http.MethodGet,
			Handler: c.MakeHandler(c.handler.GetAllOrders, nil),
		},
		kit.Endpoint{
			URI:     "/orders/batch",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.BatchBuyGas, &vo.BatchBuyRequest{}),
		},
	)
	return ends
}
