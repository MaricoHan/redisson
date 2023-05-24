package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/evm"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/types"
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/msgs"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IMsgs interface {
	GetNFTHistory(ctx context.Context, params evm.NftOperationHistoryByNftId) (*evm.NftOperationHistoryByNftIdRes, error)
	GetAccountHistory(ctx context.Context, params dto.AccountsInfo) (*dto.AccountOperationEVMRecordRes, error)
}

type msgs struct {
	logger  *log.Logger
	timeout context.Context
}

func NewMsgs(logger *log.Logger) *msgs {
	return &msgs{logger: logger}
}

func (s *msgs) GetNFTHistory(ctx context.Context, params evm.NftOperationHistoryByNftId) (*evm.NftOperationHistoryByNftIdRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "GetNFTHistory")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	sort, ok := pb.SORTS_value[params.SortBy]
	if !ok {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.NFTHistoryRequest{
		ProjectId:  params.ProjectID,
		NftId:      params.NftId,
		Signer:     params.Signer,
		TxHash:     params.TxHash,
		PageKey:    params.PageKey,
		CountTotal: params.CountTotal,
		Limit:      params.Limit,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		ClassId:    params.ClassID,
		SortBy:     pb.SORTS(sort),
	}
	req.Operation = params.Operation

	resp := &pb.NFTHistoryResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.EvmMsgsClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.NFTHistory(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &evm.NftOperationHistoryByNftIdRes{
		PageRes: dto.PageRes{
			PrevPageKey: resp.PrevPageKey,
			NextPageKey: resp.NextPageKey,
			Limit:       resp.Limit,
			TotalCount:  resp.TotalCount,
		},
		OperationRecords: []*evm.OperationRecord{},
	}
	result.TotalCount = resp.TotalCount
	var operationRecords []*evm.OperationRecord
	for _, item := range resp.Data {
		var operationRecord = &evm.OperationRecord{
			TxHash:    item.TxHash,
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

func (s *msgs) GetAccountHistory(ctx context.Context, params dto.AccountsInfo) (*dto.AccountOperationEVMRecordRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "GetAccountHistory")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	sort, ok := pb.SORTS_value[params.SortBy]
	if !ok {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.AccountHistoryRequest{
		ProjectId:  params.ProjectID,
		PageKey:    params.PageKey,
		CountTotal: params.CountTotal,
		Limit:      params.Limit,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		SortBy:     pb.SORTS(sort),
		Account:    params.Account,
		Module:     params.OperationModule,
		Operation:  params.Operation,
		TxHash:     params.TxHash,
	}

	resp := &pb.AccountHistoryResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.EvmMsgsClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.AccountHistory(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.AccountOperationEVMRecordRes{
		PageRes: dto.PageRes{
			PrevPageKey: resp.PrevPageKey,
			NextPageKey: resp.NextPageKey,
			Limit:       resp.Limit,
			TotalCount:  resp.TotalCount,
		},
		OperationRecords: []*dto.AccountOperationEVMRecords{},
	}
	var accountOperationRecords []*dto.AccountOperationEVMRecords
	for _, item := range resp.Data {
		accountOperationRecord := &dto.AccountOperationEVMRecords{
			TxHash:    item.TxHash,
			Module:    item.Module,
			Operation: item.Operation,
			Signer:    item.Signer,
			Timestamp: item.Timestamp,
		}
		if item.NftMsg != "" {
			typeJsonNft := types.JSON{}
			if err := json.Unmarshal([]byte(item.NftMsg), &typeJsonNft); err != nil {
				return nil, err
			}
			accountOperationRecord.NftMsg = &typeJsonNft
		}
		accountOperationRecords = append(accountOperationRecords, accountOperationRecord)
	}
	if accountOperationRecords != nil {
		result.OperationRecords = accountOperationRecords
	}

	return result, nil
}