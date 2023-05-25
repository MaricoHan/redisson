package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/buy"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"
)

type IBusiness interface {
	GetOrderInfo(ctx context.Context, params dto.GetOrder) (*dto.OrderInfo, error)
	GetAllOrders(ctx context.Context, params dto.GetAllOrder) (*dto.OrderOperationRes, error)
	BuildOrder(ctx context.Context, params dto.BuildOrderInfo) (*dto.BuyResponse, error)
	BatchBuyGas(ctx context.Context, params dto.BatchBuyGas) (*dto.BuyResponse, error)
}

type business struct {
	logger *log.Logger
}

func NewBusiness(logger *log.Logger) *business {
	return &business{logger: logger}
}

func (s *business) GetOrderInfo(ctx context.Context, params dto.GetOrder) (*dto.OrderInfo, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "GetOrderInfo")

	req := pb.OrderShowRequest{
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}
	resp := &pb.BuyOrderShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	// 非托管模式仅支持文昌链-天舟；托管模式仅支持文昌链-DDC
	if (params.AccessMode != entity.UNMANAGED || mapKey != constant.TianzhouEVM) && (params.AccessMode != entity.MANAGED || mapKey != constant.TianheDDC) {
		return nil, errors2.ErrNotImplemented
	}
	grpcClient, ok := initialize.BusineessClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrGrpc)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.OrderInfo{
		OperationId: resp.OperationId,
		Status:      resp.Status,
		Message:     resp.Message,
		Account:     resp.Address,
		Amount:      resp.Amount,
		Number:      resp.Number,
		CreateTime:  resp.CreatedAt,
		UpdateTime:  resp.UpdatedAt,
		OrderType:   resp.Type,
	}
	return result, nil

}

func (s *business) GetAllOrders(ctx context.Context, params dto.GetAllOrder) (*dto.OrderOperationRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "GetAllOrders")
	sorts := strings.Split(params.SortBy, "_")

	if len(sorts) != 2 {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	var sort pb.Sorts
	if sorts[0] == "DATE" {
		sort = pb.Sorts_CREATED_AT
	} else {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	var rule pb.SortRule
	if sorts[1] == "DESC" {
		rule = pb.SortRule_DESC
	} else if sorts[1] == "ASC" {
		rule = pb.SortRule_ASC
	} else {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.OrderListRequest{
		ProjectId:   params.ProjectId,
		OperationId: params.OperationId,
		PageKey:     params.PageKey,
		Limit:       params.Limit,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		SortBy:      sort,
		SortRule:    rule,
		Address:     params.Account,
		CountTotal:  params.CountTotal,
		Status:      params.Status,
	}

	resp := &pb.BuyOrderListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)

	// 非托管模式仅支持文昌链-天舟；托管模式仅支持文昌链-DDC
	if (params.AccessMode != entity.UNMANAGED || params.Code != constant.Tianzhou) && (params.AccessMode != entity.MANAGED || mapKey != constant.TianheDDC) {
		return nil, errors2.ErrNotImplemented
	}

	grpcClient, ok := initialize.BusineessClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.List(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.OrderOperationRes{
		PageRes: dto.PageRes{
			PrevPageKey: resp.PrevPageKey,
			NextPageKey: resp.NextPageKey,
			Limit:       uint32(resp.Limit),
			TotalCount:  resp.TotalCount,
		},
		OrderInfos: []*dto.OrderInfo{},
	}

	var orderOperationRes []*dto.OrderInfo
	for _, item := range resp.Data {
		var orderInfo = &dto.OrderInfo{
			OperationId: item.OperationId,
			Status:      item.Status,
			Message:     item.Message,
			Account:     item.Address,
			Amount:      item.Amount,
			Number:      item.Number,
			CreateTime:  item.CreatedAt,
			UpdateTime:  item.UpdatedAt,
			OrderType:   item.Type,
		}
		orderOperationRes = append(orderOperationRes, orderInfo)
	}
	if orderOperationRes != nil {
		result.OrderInfos = orderOperationRes
	}

	return result, nil

}

func (s *business) BuildOrder(ctx context.Context, params dto.BuildOrderInfo) (*dto.BuyResponse, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "BuildOrder")

	req := pb.BuyRequest{
		ProjectId:   params.ProjectID,
		Address:     params.Address,
		Amount:      params.Amount,
		OperationId: params.OperationId,
	}
	resp := &pb.BuyResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	// 非托管模式仅支持文昌链-天舟充值 gas；托管模式仅支持文昌链-DDC
	if params.OrderType != constant.OrderTypeGas || params.AccessMode != entity.UNMANAGED || mapKey != constant.TianzhouEVM {
		if params.OrderType != constant.OrderTypeGas && params.AccessMode == entity.UNMANAGED && mapKey == constant.TianzhouEVM {
			return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderType)
		}
		return nil, errors2.ErrNotImplemented
	}
	grpcClient, ok := initialize.BusineessClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	switch params.OrderType {
	case constant.OrderTypeGas:
		resp, err = grpcClient.BuyGas(ctx, &req)
		if err != nil {
			logger.Error("request err:", err.Error())
			return nil, err
		}
	default:
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderType)
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.BuyResponse{}
	return result, nil

}

func (s *business) BatchBuyGas(ctx context.Context, params dto.BatchBuyGas) (*dto.BuyResponse, error) {
	logger := s.logger.WithContext(ctx).WithFields(map[string]interface{}{
		"func":   "BatchBuyGas",
		"params": params,
	})

	req := pb.BatchBuyRequest{
		ProjectId:   params.ProjectID,
		List:        params.List,
		OperationId: params.OperationId,
	}
	resp := &pb.BatchBuyResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	// 非托管模式仅支持文昌链-天舟；
	if params.AccessMode != entity.UNMANAGED || mapKey != constant.TianzhouEVM {
		return nil, errors2.ErrNotImplemented
	}
	grpcClient, ok := initialize.BusineessClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchBuyGas(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.BuyResponse{}
	return result, nil
}
