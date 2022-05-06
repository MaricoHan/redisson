package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type EmptionController struct {
	BaseController
	handler handlers.IBusiness
}

func NewEmptionController(bc BaseController, handler handlers.IBusiness) kit.IController {
	return EmptionController{bc, handler}
}

func (c EmptionController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/orders/{order_id}",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.GetOrderInfo, nil),
		},
		kit.Endpoint{
			URI:     "/orders",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.BuildOrder, &vo.BuyRequest{}),
		},
		kit.Endpoint{
			URI:     "/orders",
			Method:  http.MethodGet,
			Handler: c.makeHandler(c.handler.GetAllOrders, nil),
		},
	)
	return ends
}
