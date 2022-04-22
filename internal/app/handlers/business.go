package handlers

import (
	"context"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"regexp"
)

type IBusiness interface {
	GetOrderInfo(ctx context.Context, _ interface{}) (interface{}, error)
	BuildOrder(ctx context.Context, _ interface{}) (interface{}, error)
	GetAllOrders(ctx context.Context, _ interface{}) (interface{}, error)
}

type Business struct {
	base
	pageBasic
	svc services.IBusiness
}

func NewBusiness(svc services.IBusiness) *Business {
	return &Business{svc: svc}
}

func (h *Business) GetOrderInfo(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)

	params := dto.GetOrder{
		ProjectID: authData.ProjectId,
		Module:    authData.Module,
		OrderId:   h.GetOrderId(ctx),
		Code:      authData.Code,
	}
	return h.svc.GetOrderInfo(params)
}

func (h *Business) GetAllOrders(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)

	params := dto.GetAllOrder{
		Page:      dto.Page{},
		Module:    authData.Module,
		ProjectId: authData.ProjectId,
		Account:   h.GetAccount(ctx),
		Code:      authData.Code,
		Status:    h.GetStatus(ctx),
	}
	offset, err := h.Offset(ctx)
	if err != nil {
		return nil, err
	}
	params.Offset = offset

	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, err
	}
	params.Limit = limit

	if params.Limit == 0 {
		params.Limit = 10
	}

	params.StartDate = h.StartDate(ctx)
	params.EndDate = h.EndDate(ctx)
	params.SortBy = h.SortBy(ctx)
	return h.svc.GetAllOrders(params)
}

func (h *Business) BuildOrder(ctx context.Context, request interface{}) (interface{}, error) {
	OrderRes := request.(*vo.BuyRequest)
	authData := h.AuthData(ctx)

	if len(OrderRes.OrderId) == 0 {
		return nil, errors2.New(errors2.ClientParams, "order_id is a required field")
	}
	if len(OrderRes.OrderType) == 0 {
		return nil, errors2.New(errors2.ClientParams, "order_type is a required field")
	}

	orderId := OrderRes.OrderId

	if len(orderId) < 10 || len(orderId) > 36 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderIDLen)
	}
	ok, err := regexp.MatchString("^([A-Za-z0-9_]){10,36}$", orderId)
	if !ok || err != nil {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderID)
	}

	ok, err = regexp.MatchString("([A-Za-z])+", orderId)
	if !ok || err != nil {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderID)
	}
	ok, err = regexp.MatchString("([0-9])+", orderId)
	if !ok || err != nil {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderID)
	}
	ok, err = regexp.MatchString("([_])+", orderId)
	if !ok || err != nil {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderID)
	}

	if OrderRes.Amount < 100 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderAmount)
	}
	if OrderRes.Amount%100 != 0 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrAmountFormat)
	}

	params := dto.BuildOrderInfo{
		ProjectID: authData.ProjectId,
		ChainId:   authData.ChainId,
		Address:   OrderRes.Account,
		Amount:    OrderRes.Amount,
		Module:    authData.Module,
		OrderType: OrderRes.OrderType,
		OrderId:   OrderRes.OrderId,
		Code:      authData.Code,
	}
	return h.svc.BuildOrder(params)
}

func (h *Business) GetAddress(ctx context.Context) string {
	address := ctx.Value("address")
	if address == nil {
		return ""
	}
	return address.(string)
}

func (h *Business) GetAccount(ctx context.Context) string {
	account := ctx.Value("account")
	if account == nil {
		return ""
	}
	return account.(string)
}

func (h *Business) GetOrderType(ctx context.Context) string {
	orderType := ctx.Value("order_type")
	if orderType == nil {
		return ""
	}
	return orderType.(string)
}

func (h *Business) GetOrderId(ctx context.Context) string {
	orderId := ctx.Value("order_id")
	if orderId == nil {
		return ""
	}
	return orderId.(string)
}

func (h *Business) GetStatus(ctx context.Context) string {
	status := ctx.Value("status")
	if status == nil {
		return ""
	}
	return status.(string)
}
