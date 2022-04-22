package services

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/types"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/msgs"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"time"
)

type IMsgs interface {
	GetNFTHistory(params dto.NftOperationHistoryByNftId) (*dto.NftOperationHistoryByNftIdRes, error)
	GetAccountHistory(params dto.AccountsInfo) (*dto.AccountOperationRecordRes, error)
}

type msgs struct {
	logger  *log.Logger
	timeout context.Context
}

func NewMsgs(logger *log.Logger) *msgs {
	return &msgs{logger: logger}
}

func (s *msgs) GetNFTHistory(params dto.NftOperationHistoryByNftId) (*dto.NftOperationHistoryByNftIdRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "msgs"
	logFields["func"] = "GetNFTHistory"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	sort, ok := pb.Sorts_value[params.SortBy]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.NFTHistoryRequest{
		ProjectId: params.ProjectID,
		NftId:     params.NftId,
		Signer:    params.Signer,
		TxHash:    params.Txhash,
		Offset:    params.Offset,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		ClassId:   params.ClassID,
		SortBy:    pb.Sorts(sort),
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

	resp := &pb.NFTHistoryResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MsgsClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.NFTHistory(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftOperationHistoryByNftIdRes{
		PageRes: dto.PageRes{
			Offset: resp.Offset,
			Limit:  resp.Limit,
		},
		OperationRecords: []*dto.OperationRecord{},
	}
	result.TotalCount = resp.TotalCount
	var operationRecords []*dto.OperationRecord
	for _, item := range resp.Data {
		var operationRecord = &dto.OperationRecord{
			Txhash:    item.TxHash,
			Operation: item.Operation,
			Signer:    item.Signer,
			Recipient: item.Recipient,
			Timestamp: item.Timestamp,
		}
		operationRecords = append(operationRecords, operationRecord)
	}
	if operationRecords != nil {
		result.OperationRecords = operationRecords
	}

	return result, nil
}

func (s *msgs) GetAccountHistory(params dto.AccountsInfo) (*dto.AccountOperationRecordRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "account"
	logFields["func"] = "GetAccountHistory"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	var operation, module int32
	var ok bool
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	if params.OperationModule == "" {
		operation = 0
		module = 0
	}

	if params.OperationModule != "" {
		module, ok = pb.Module_value[params.OperationModule]
		if !ok {
			log.WithFields(logFields).Error(errors2.ErrModule)
			return nil, errors2.New(errors2.ClientParams, errors2.ErrModule)
		}
		if params.Operation == "" {

			operation = 0
		} else {
			operation, ok = pb.Operation_value[params.Operation]
			if !ok {
				log.WithFields(logFields).Error(errors2.ErrOperation)
				return nil, errors2.New(errors2.ClientParams, errors2.ErrOperation)
			}
		}

	}

	sort, ok := pb.Sorts_value[params.SortBy]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.AccountHistoryRequest{
		ProjectId: params.ProjectID,
		Offset:    params.Offset,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		SortBy:    pb.Sorts(sort),
		Address:   params.Account,
		Module:    pb.Module(module),
		Operation: pb.Operation(operation),
	}

	resp := &pb.AccountHistoryResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MsgsClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.AccountHistory(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.AccountOperationRecordRes{
		PageRes: dto.PageRes{
			Offset:     resp.Offset,
			Limit:      resp.Limit,
			TotalCount: resp.TotalCount,
		},
		OperationRecords: []*dto.AccountOperationRecords{},
	}
	var accountOperationRecords []*dto.AccountOperationRecords
	for _, item := range resp.Data {
		typeJson := types.JSON{}
		err := json.Unmarshal([]byte(item.Message), &typeJson)
		if err != nil {
			return nil, err
		}
		accountOperationRecord := &dto.AccountOperationRecords{
			TxHash:    item.TxHash,
			Module:    item.Module,
			Operation: item.Operation,
			Signer:    item.Signer,
			Timestamp: item.Timestamp,
			Message:   typeJson,
		}
		accountOperationRecords = append(accountOperationRecords, accountOperationRecord)
	}
	if accountOperationRecords != nil {
		result.OperationRecords = accountOperationRecords
	}

	return result, nil
}
