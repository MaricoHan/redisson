package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"
	"gitlab.bianjie.ai/avata/utils/errors/v2/common"
)

type IBusiness interface {
	GetOrderInfo(ctx context.Context, _ interface{}) (interface{}, error)
	BuildOrder(ctx context.Context, _ interface{}) (interface{}, error)
	GetAllOrders(ctx context.Context, _ interface{}) (interface{}, error)
	BatchBuyGas(ctx context.Context, _ interface{}) (interface{}, error)
}

type Business struct {
	base.Base
	base.PageBasic
	svc services.IBusiness
}

func NewBusiness(svc services.IBusiness) *Business {
	return &Business{svc: svc}
}

func (h *Business) GetOrderInfo(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)

	params := dto.GetOrder{
		ProjectID:   authData.ProjectId,
		Module:      authData.Module,
		OperationId: h.GetOperationId(ctx),
		Code:        authData.Code,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.GetOrderInfo(ctx, params)
}

func (h *Business) GetAllOrders(ctx context.Context, _ interface{}) (interface{}, error) {
	authData := h.AuthData(ctx)

	params := dto.GetAllOrder{
		Module:     authData.Module,
		ProjectId:  authData.ProjectId,
		Account:    h.GetAccount(ctx),
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}

	status, err := h.GetStatus(ctx)
	if err != nil {
		return nil, err
	}
	params.Status = uint32(status)

	params.PageKey = h.PageKey(ctx)
	countTotal, err := h.CountTotal(ctx)
	if err != nil {
		return nil, errors2.New(errors2.ClientParams, fmt.Sprintf(common.ERR_INVALID_VALUE, "count_total"))
	}
	params.CountTotal = countTotal

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
	return h.svc.GetAllOrders(ctx, params)
}

func (h *Business) BuildOrder(ctx context.Context, request interface{}) (interface{}, error) {
	OrderRes := request.(*vo.BuyRequest)
	authData := h.AuthData(ctx)

	if len(OrderRes.OperationId) == 0 {
		return nil, errors2.New(errors2.ClientParams, "operation_id is a required field")
	}
	if OrderRes.OrderType == 0 {
		return nil, errors2.New(errors2.ClientParams, "order_type is a required field")
	}

	operationId := strings.TrimSpace(OrderRes.OperationId)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	if OrderRes.Amount < 100 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderAmount)
	}
	if OrderRes.Amount%100 != 0 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrAmountFormat)
	}

	params := dto.BuildOrderInfo{
		ProjectID:   authData.ProjectId,
		ChainId:     authData.ChainId,
		Address:     OrderRes.Account,
		Amount:      OrderRes.Amount,
		Module:      authData.Module,
		OrderType:   OrderRes.OrderType,
		OperationId: OrderRes.OperationId,
		Code:        authData.Code,
		AccessMode:  authData.AccessMode,
	}
	return h.svc.BuildOrder(ctx, params)
}

func (h *Business) BatchBuyGas(ctx context.Context, request interface{}) (interface{}, error) {
	OrderRes := request.(*vo.BatchBuyRequest)
	authData := h.AuthData(ctx)

	if len(OrderRes.OperationId) == 0 {
		return nil, errors2.New(errors2.ClientParams, "operation_id is a required field")
	}

	if len(OrderRes.List) == 0 {
		return nil, errors2.New(errors2.ClientParams, "list is a required field")
	}

	operationId := strings.TrimSpace(OrderRes.OperationId)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	params := dto.BatchBuyGas{
		ProjectID:   authData.ProjectId,
		ChainId:     authData.ChainId,
		Module:      authData.Module,
		List:        OrderRes.List,
		OperationId: OrderRes.OperationId,
		Code:        authData.Code,
		AccessMode:  authData.AccessMode,
	}
	return h.svc.BatchBuyGas(ctx, params)
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

func (h *Business) GetOperationId(ctx context.Context) string {
	operationId := ctx.Value("operation_id")
	if operationId == nil {
		return ""
	}
	return operationId.(string)
}

func (h *Business) GetStatus(ctx context.Context) (int64, error) {
	value := ctx.Value("status")
	if value == nil {
		return 0, nil
	}

	status, err := strconv.ParseInt(value.(string), 10, 64)
	if err != nil || status < 0 || status > 3 {
		return 0, errors2.New(errors2.ClientParams, fmt.Sprintf(common.ERR_INVALID_VALUE, "status"))
	}
	return status, nil
}
