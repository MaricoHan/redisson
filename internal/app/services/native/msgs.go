package native

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/volatiletech/sqlboiler/types"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/msgs"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"
)

type IMsgs interface {
	GetNFTHistory(ctx context.Context, params nft.NftOperationHistoryByNftId) (*nft.NftOperationHistoryByNftIdRes, error)
	GetAccountHistory(ctx context.Context, params dto.AccountsInfo) (*dto.AccountOperationRecordRes, error)
	GetMTHistory(ctx context.Context, params mt.MTOperationHistoryByMTId) (*mt.MTOperationHistoryByMTIdRes, error)
}

type msgs struct {
	logger  *log.Logger
	timeout context.Context
}

func NewMsgs(logger *log.Logger) *msgs {
	return &msgs{logger: logger}
}

func (s *msgs) GetNFTHistory(ctx context.Context, params nft.NftOperationHistoryByNftId) (*nft.NftOperationHistoryByNftIdRes, error) {
	logger := s.logger.WithField("params", params).WithField("func", "GetNFTHistory")

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
		ProjectId: params.ProjectID,
		NftId:     params.NftId,
		Signer:    params.Signer,
		TxHash:    params.Txhash,
		PageKey:   params.PageKey,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		ClassId:   params.ClassID,
		SortBy:    pb.SORTS(sort),
		Operation: params.Operation,
	}
	resp := &pb.NFTHistoryResponse{}

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMsgClientMap[mapKey]
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
	result := &nft.NftOperationHistoryByNftIdRes{
		PageRes: dto.PageRes{
			PrevPageKey: resp.PrevPageKey,
			NextPageKey: resp.NextPageKey,
			Limit:       resp.Limit,
		},
		OperationRecords: []*nft.OperationRecord{},
	}
	result.TotalCount = resp.TotalCount
	var operationRecords []*nft.OperationRecord
	for _, item := range resp.Data {
		var operationRecord = &nft.OperationRecord{
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

func (s *msgs) GetAccountHistory(ctx context.Context, params dto.AccountsInfo) (*dto.AccountOperationRecordRes, error) {
	logger := s.logger.WithField("params", params).WithField("func", "GetAccountHistory")

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
		Limit:      params.Limit,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		SortBy:     pb.SORTS(sort),
		Address:    params.Account,
		Module:     params.OperationModule,
		Operation:  params.Operation,
		TxHash:     params.TxHash,
		CountTotal: params.CountTotal,
	}

	resp := &pb.AccountHistoryResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMsgClientMap[mapKey]
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
	result := &dto.AccountOperationRecordRes{
		PageRes: dto.PageRes{
			NextPageKey: resp.NextPageKey,
			PrevPageKey: resp.PrevPageKey,
			Limit:       resp.Limit,
			TotalCount:  resp.TotalCount,
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

func (s *msgs) GetMTHistory(ctx context.Context, params mt.MTOperationHistoryByMTId) (*mt.MTOperationHistoryByMTIdRes, error) {
	logger := s.logger.WithField("params", params).WithField("func", "GetMTHistory")

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
	req := pb.MTHistoryRequest{
		ProjectId: params.ProjectID,
		PageKey:   params.PageKey,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		SortBy:    pb.SORTS(sort),
		Signer:    params.Signer,
		TxHash:    params.Txhash,
		MtId:      params.MTId,
		ClassId:   params.ClassID,
	}

	req.Operation = params.Operation

	resp := &pb.MTHistoryResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMsgClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.MTHistory(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &mt.MTOperationHistoryByMTIdRes{
		PageRes: dto.PageRes{
			NextPageKey: resp.NextPageKey,
			PrevPageKey: resp.PrevPageKey,
			TotalCount:  resp.TotalCount,
			Limit:       resp.Limit,
		},
		OperationRecords: []*mt.MTOperationRecord{},
	}
	var operationRecords []*mt.MTOperationRecord
	for _, item := range resp.Data {
		var operationRecord = &mt.MTOperationRecord{
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
