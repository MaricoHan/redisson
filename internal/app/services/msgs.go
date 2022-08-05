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
	GetMTHistory(params dto.MTOperationHistoryByMTId) (*dto.MTOperationHistoryByMTIdRes, error)
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

	sort, ok := pb.SORTS_value[params.SortBy]
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
		SortBy:    pb.SORTS(sort),
	}
	if params.Operation != "" {
		operation, ok := pb.OPERATION_value[params.Operation]
		if !ok {
			if !ok {
				log.WithFields(logFields).Error(errors2.ErrOperation)
				return nil, errors2.New(errors2.ClientParams, errors2.ErrOperation)
			}
		}
		req.Operation = pb.OPERATION(operation)
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
			Operation: item.Operation.String(),
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
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	var operation, module int32
	var ok bool
	if params.OperationModule == "" {
		operation = 0
		module = 0
	} else {
		module, ok = pb.Module_value[params.OperationModule]
		if !ok {
			log.WithFields(logFields).Error(errors2.ErrModule)
			return nil, errors2.New(errors2.ClientParams, errors2.ErrModule)
		}

		if params.Operation == "" {
			operation = 0
		} else if operation, ok = pb.OPERATION_value[params.Operation]; !ok {
			log.WithFields(logFields).Error(errors2.ErrOperation)
			return nil, errors2.New(errors2.ClientParams, errors2.ErrOperation)
		}
	}

	sort, ok := pb.SORTS_value[params.SortBy]
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
		SortBy:    pb.SORTS(sort),
		Address:   params.Account,
		Module:    pb.Module(module),
		Operation: pb.OPERATION(operation),
		TxHash:    params.TxHash,
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
			TxHash:      item.TxHash,
			Module:      item.Module.String(),
			Operation:   item.Operation.String(),
			Signer:      item.Signer,
			Timestamp:   item.Timestamp,
			Message:     &typeJson,
			GasFee:      item.GasFee,
			BusinessFee: item.BusinessFee,
		}
		if item.NftMsg != "" {
			typeJsonNft := types.JSON{}
			if err := json.Unmarshal([]byte(item.NftMsg), &typeJsonNft); err != nil {
				return nil, err
			}
			accountOperationRecord.NftMsg = &typeJsonNft
		}
		if item.MtMsg != "" {
			typeJsonMt := types.JSON{}
			if err := json.Unmarshal([]byte(item.MtMsg), &typeJsonMt); err != nil {
				return nil, err
			}
			accountOperationRecord.MtMsg = &typeJsonMt
		}
		accountOperationRecords = append(accountOperationRecords, accountOperationRecord)
	}
	if accountOperationRecords != nil {
		result.OperationRecords = accountOperationRecords
	}

	return result, nil
}

func (m *msgs) GetMTHistory(params dto.MTOperationHistoryByMTId) (*dto.MTOperationHistoryByMTIdRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "msgs"
	logFields["func"] = "GetMTHistory"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	sort, ok := pb.SORTS_value[params.SortBy]
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
		SortBy:    pb.SORTS(sort),
		Signer:    params.Signer,
		TxHash:    params.Txhash,
		MtId:      params.MTId,
		ClassId:   params.ClassID,
	}
	if params.Operation != "" {
		operation, ok := pb.OPERATION_value[params.Operation]
		if !ok {
			if !ok {
				log.WithFields(logFields).Error(errors2.ErrOperation)
				return nil, errors2.New(errors2.ClientParams, errors2.ErrOperation)
			}
		}
		req.Operation = pb.OPERATION(operation)
	}

	resp := &pb.MTHistoryResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MsgsClientMap[mapKey]
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
			Operation: item.Operation.String(),
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
