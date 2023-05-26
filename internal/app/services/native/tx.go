package native

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/types"
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/tx"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"

	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
)

type ITx interface {
	TxResult(ctx context.Context, params dto.TxResultByTxHash) (*dto.TxResultRes, error)
	//TxQueueInfo(ctx context.Context, params dto.TxQueueInfo) (*dto.TxQueueInfoRes, error)
}

type tx struct {
	logger *log.Entry
}

func NewTx(logger *log.Logger) *tx {
	return &tx{logger: logger.WithField("model", "tx")}
}

func (t *tx) TxResult(ctx context.Context, params dto.TxResultByTxHash) (*dto.TxResultRes, error) {
	logger := t.logger.WithContext(ctx).WithField("params", params).WithField("func", "TxResult")
	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.TxShowRequest{
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}
	resp := &pb.TxShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeTxClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		logger.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil || resp.Detail == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := new(dto.TxResultRes)
	status := resp.Detail.Status
	result.Module = resp.Detail.Module
	result.Operation = resp.Detail.Operation
	result.TxHash = ""
	result.Status = uint32(status)
	if status == pb.STATUS_success || status == pb.STATUS_failed {
		result.TxHash = resp.Detail.Hash
	}

	result.Message = resp.Detail.Message
	result.BlockHeight = resp.Detail.BlockHeight
	result.Timestamp = resp.Detail.Timestamp

	if resp.Detail.Nft != "" {
		result.Nft = new(types.JSON)
		err = json.Unmarshal([]byte(resp.Detail.Nft), &result.Nft)
		if err != nil {
			logger.WithError(err).Error("nft unmarshal failed")
			return nil, errors2.ErrInternal
		}
	}

	if resp.Detail.Mt != "" {
		result.Mt = new(types.JSON)
		err = json.Unmarshal([]byte(resp.Detail.Mt), &result.Mt)
		if err != nil {
			logger.WithError(err).Error("mt unmarshal failed")
			return nil, errors2.ErrInternal
		}
	}

	if resp.Detail.Record != new(pb.Record) {
		result.Record = resp.Detail.Record
	}
	return result, nil
}

//func (t *tx) TxQueueInfo(ctx context.Context, params dto.TxQueueInfo) (*dto.TxQueueInfoRes, error) {
//	logger := t.logger.WithContext(ctx).WithField("params", params).WithField("func", "TxQueueInfo")
//
//	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
//	defer cancel()
//	req := pb_queue.TxQueueShowRequest{
//		ProjectId:   params.ProjectID,
//		OperationId: params.OperationId,
//		Code:        params.Code,
//		Module:      params.Module,
//	}
//	resp := &pb_queue.TxQueueShowResponse{}
//	var err error
//	resp, err = initialize.TxQueueClient.Show(ctx, &req)
//	if err != nil {
//		logger.WithError(err).Error("request err")
//		return nil, err
//	}
//	if resp == nil {
//		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
//	}
//	result := new(dto.TxQueueInfoRes)
//	result.QueueTotal = resp.QueueTotal
//	result.QueueRequestTime = resp.QueueRequestTime
//	result.QueueCostTime = resp.QueueCostTime
//	result.TxQueuePosition = resp.TxQueuePosition
//	result.TxRequestTime = resp.TxRequestTime
//	result.TxCostTime = resp.TxCostTime
//	result.TxMessage = resp.TxMessage
//
//	return result, nil
//}
