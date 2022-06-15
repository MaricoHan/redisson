package services

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/mt_msgs"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IMTMsgs interface {
	GetMTHistory(params dto.MTOperationHistoryByMTId) (*dto.MTOperationHistoryByMTIdRes, error)
}

type mtMsgs struct {
	logger  *log.Logger
	timeout context.Context
}

func NewMTMsgs(logger *log.Logger) *mtMsgs {
	return &mtMsgs{logger: logger}
}

func (m mtMsgs) GetMTHistory(params dto.MTOperationHistoryByMTId) (*dto.MTOperationHistoryByMTIdRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt_msgs"
	logFields["func"] = "GetMTHistory"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	sort, ok := pb.Sorts_value[params.SortBy]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.MTHistoryRequest{
		ProjectId: params.ProjectID,
		Offset:    params.Offset,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		SortBy:    pb.Sorts(sort),
		Signer:    params.Signer,
		TxHash:    params.Txhash,
		MtId:      params.MTId,
		ClassId:   params.ClassID,
	}
	if params.Operation != "" {
		operation, ok := pb.Operation_value[params.Operation]
		if !ok {
			if !ok {
				log.WithFields(logFields).Error(errors2.ErrOperation)
				return nil, errors2.New(errors2.ClientParams, errors2.ErrOperation)
			}
		}
		req.Operation = pb.Operation(operation)
	}

	resp := &pb.MTHistoryResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTMsgsClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.MTHistory(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.MTOperationHistoryByMTIdRes{
		PageRes: dto.PageRes{
			Offset: resp.Offset,
			Limit:  resp.Limit,
		},
		OperationRecords: []*dto.MTOperationRecord{},
	}
	result.TotalCount = resp.TotalCount
	var operationRecords []*dto.MTOperationRecord
	for _, item := range resp.Data {
		var operationRecord = &dto.MTOperationRecord{
			Txhash:    item.TxHash,
			Operation: item.Operation,
			Signer:    item.Signer,
			Recipient: item.Recipient,
			Amount:    item.Amount,
			Timestamp: item.Timestamp,
		}
		operationRecords = append(operationRecords, operationRecord)
	}
	if operationRecords != nil {
		result.OperationRecords = operationRecords
	}

	return result, nil
}
