package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/buy"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"strings"
	"time"
)

type IBusiness interface {
	GetOrderInfo(params dto.GetOrder) (*dto.OrderInfo, error)
	GetAllOrders(params dto.GetAllOrder) (*dto.OrderOperationRes, error)
	BuildOrder(params dto.BuildOrderInfo) (*dto.BuyResponse, error)
}

type business struct {
	logger *log.Logger
}

func NewBusiness(logger *log.Logger) *business {
	return &business{logger: logger}
}

func (s *business) GetOrderInfo(params dto.GetOrder) (*dto.OrderInfo, error) {
	logFields := log.Fields{}
	logFields["model"] = "order"
	logFields["func"] = "GetOrderInfo"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.OrderShowRequest{
		ProjectId: params.ProjectID,
		OrderId:   params.OrderId,
	}
	resp := &pb.BuyOrderShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.BusineessClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrGrpc)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.OrderInfo{
		OrderId:    resp.OrderId,
		Status:     strings.ToLower(pb.Status_name[int32(resp.Status)]),
		Message:    resp.Message,
		Account:    resp.Address,
		Amount:     resp.Amount,
		Number:     resp.Number,
		CreateTime: resp.CreatedAt,
		UpdateTime: resp.UpdatedAt,
		OrderType:  resp.Type,
	}
	return result, nil

}

func (s *business) GetAllOrders(params dto.GetAllOrder) (*dto.OrderOperationRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "order"
	logFields["func"] = "GetAllOrders"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	sorts := strings.Split(params.SortBy, "_")

	if len(sorts) != 2 {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	var sort pb.Sorts
	if sorts[0] == "DATE" {
		sort = pb.Sorts_ID
	} else {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	var rule pb.SortRule
	if sorts[1] == "DESC" {
		rule = pb.SortRule_DESC
	} else if sorts[1] == "ASC" {
		rule = pb.SortRule_ASC
	} else {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.OrderListRequest{
		ProjectId: params.ProjectId,
		OrderId:   params.OrderId,
		Offset:    params.Offset,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		SortBy:    sort,
		SortRule:  rule,
		Address:   params.Account,
		//Status: pb.Status(status),

	}
	if params.Status != "" {
		status, ok := pb.Status_value[strings.ToUpper(params.Status)]
		if !ok {
			log.WithFields(logFields).Error(errors2.ErrStatus)
			return nil, errors2.New(errors2.ClientParams, errors2.ErrStatus)
		}
		req.Status = pb.Status(status)
	}

	resp := &pb.BuyOrderListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)

	grpcClient, ok := initialize.BusineessClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.List(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.OrderOperationRes{
		PageRes: dto.PageRes{
			Offset:     resp.Offset,
			Limit:      resp.Limit,
			TotalCount: resp.TotalCount,
		},
		OrderInfos: []*dto.OrderInfo{},
	}

	var orderOperationRes []*dto.OrderInfo
	for _, item := range resp.Data {
		var orderInfo = &dto.OrderInfo{
			OrderId:    item.OrderId,
			Status:     strings.ToLower(pb.Status_name[int32(item.Status)]),
			Message:    item.Message,
			Account:    item.Address,
			Amount:     item.Amount,
			Number:     item.Number,
			CreateTime: item.CreatedAt,
			UpdateTime: item.UpdatedAt,
			OrderType:  item.Type,
		}
		orderOperationRes = append(orderOperationRes, orderInfo)
	}
	if orderOperationRes != nil {
		result.OrderInfos = orderOperationRes
	}

	return result, nil

}

func (s *business) BuildOrder(params dto.BuildOrderInfo) (*dto.BuyResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "order"
	logFields["func"] = "BuildOrder"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.BuyRequest{
		ProjectId: params.ProjectID,
		Address:   params.Address,
		Amount:    params.Amount,
		OrderId:   params.OrderId,
	}
	resp := &pb.BuyResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.BusineessClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	switch params.OrderType {
	case constant.OrderTypeGas:
		resp, err = grpcClient.BuyGas(ctx, &req)
		if err != nil {
			log.WithFields(logFields).Error("request err:", err.Error())
			return nil, err
		}
	case constant.OrderTypeBusiness:
		resp, err = grpcClient.BuyBusiness(ctx, &req)
		if err != nil {
			log.WithFields(logFields).Error("request err:", err.Error())
			return nil, err
		}
	default:
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOrderType)
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.BuyResponse{
		OrderId: resp.OrderId,
	}
	return result, nil

}