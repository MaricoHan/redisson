package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type RightsController struct {
	BaseController
	handler handlers.IRights
}

func NewRightsController(bc BaseController, handler handlers.IRights) kit.IController {
	return RightsController{bc, handler}
}

func (r RightsController) GetEndpoints() []kit.Endpoint {
	return []kit.Endpoint{
		{
			URI:     "/rights/register",
			Method:  http.MethodPost,
			Handler: r.makeHandler(r.handler.Register, &vo.RegisterRequest{}),
		},
		{
			URI:     "/rights/register/{operation_id}",
			Method:  http.MethodPatch,
			Handler: r.makeHandler(r.handler.EditRegister, &vo.EditRegisterRequest{}),
		},
		{
			URI:     "/rights/register/{operation_id}",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.QueryRegister, nil),
		},
		{
			URI:     "/rights/user/auth",
			Method:  http.MethodPost,
			Handler: r.makeHandler(r.handler.UserAuth, &vo.UserAuthRequest{}),
		},
		{
			URI:     "/rights/user/auth/{operation_id}",
			Method:  http.MethodPatch,
			Handler: r.makeHandler(r.handler.EditUserAuth, &vo.EditUserAuthRequest{}),
		},
		{
			URI:     "/rights/user/auth",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.QueryUserAuth, nil),
		},
		{
			URI:     "/rights/dict",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.Dict, nil),
		},
		{
			URI:     "/rights/region",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.Region, nil),
		},

		{
			URI:     "/rights/delivery",
			Method:  http.MethodPost,
			Handler: r.makeHandler(r.handler.Delivery, &vo.DeliveryRequest{}),
		}, {
			URI:     "/rights/delivery/{operation_id}",
			Method:  http.MethodPatch,
			Handler: r.makeHandler(r.handler.EditDelivery, &vo.EditDeliveryRequest{}),
		}, {
			URI:     "/rights/delivery",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.DeliveryInfo, nil),
		},

		{
			URI:     "/rights/change",
			Method:  http.MethodPost,
			Handler: r.makeHandler(r.handler.Change, &vo.ChangeRequest{}),
		}, {
			URI:     "/rights/change/{operation_id}",
			Method:  http.MethodPatch,
			Handler: r.makeHandler(r.handler.EditChange, &vo.EditChangeRequest{}),
		}, {
			URI:     "/rights/change/{operation_id}",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.ChangeInfo, nil),
		},

		{
			URI:     "/rights/transfer",
			Method:  http.MethodPost,
			Handler: r.makeHandler(r.handler.Transfer, &vo.TransferRequest{}),
		}, {
			URI:     "/rights/transfer/{operation_id}",
			Method:  http.MethodPatch,
			Handler: r.makeHandler(r.handler.EditTransfer, &vo.EditTransferRequest{}),
		}, {
			URI:     "/rights/transfer/{operation_id}",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.TransferInfo, nil),
		},

		{
			URI:     "/rights/revoke",
			Method:  http.MethodPost,
			Handler: r.makeHandler(r.handler.Revoke, &vo.RevokeRequest{}),
		}, {
			URI:     "/rights/revoke/{operation_id}",
			Method:  http.MethodPatch,
			Handler: r.makeHandler(r.handler.EditRevoke, &vo.EditRevokeRequest{}),
		}, {
			URI:     "/rights/revoke/{operation_id}",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.RevokeInfo, nil),
		},

		{
			URI:     "/rights/product/{product_id}",
			Method:  http.MethodGet,
			Handler: r.makeHandler(r.handler.ProductInfo, nil),
		},
	}
}
