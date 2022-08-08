package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/tx"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
)

type ITx interface {
	TxResultByTxHash(params dto.TxResultByTxHash) (*dto.TxResultByTxHashRes, error)
}

type tx struct {
	logger *log.Entry
}

func NewTx(logger *log.Logger) *tx {
	return &tx{logger: logger.WithField("model", "tx")}
}

func (t *tx) TxResultByTxHash(params dto.TxResultByTxHash) (*dto.TxResultByTxHashRes, error) {
	log := t.logger.WithFields(map[string]interface{}{
		"func":   "TxResultByTxHash",
		"module": params.Module,
		"code":   params.Code,
	})

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.TxShowRequest{
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}
	resp := &pb.TxShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.TxClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		log.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil || resp.Detail == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := new(dto.TxResultByTxHashRes)
	status := resp.Detail.Status
	result.Module = resp.Detail.Module
	result.Type = resp.Detail.Operation
	result.TxHash = ""
	result.Status = int32(status)
	if status == pb.STATUS_success || status == pb.STATUS_failed {
		result.TxHash = resp.Detail.Hash
	}
	if resp.Detail.Tag != "" {
		var tagInterface interface{}
		err = json.Unmarshal([]byte(resp.Detail.Tag), &tagInterface)
		if err != nil {
			log.WithError(err).Error("Unmarshal failed")
			return nil, errors2.ErrInternal
		}
		result.Tag = tagInterface.(map[string]interface{})
	}

	result.Message = resp.Detail.ErrMsg
	result.BlockHeight = resp.Detail.BlockHeight
	result.Timestamp = resp.Detail.Timestamp

	if resp.Detail.Nft != "" {
		err = json.Unmarshal([]byte(resp.Detail.Nft), &result.Nft)
		if err != nil {
			log.WithError(err).Error("Unmarshal failed")
			return nil, errors2.ErrInternal
		}
		result.NftID = resp.Detail.NftId
		result.ClassID = resp.Detail.ClassId
	}
	if resp.Detail.Mt != "" {
		err = json.Unmarshal([]byte(resp.Detail.Mt), result.Mt)
		if err != nil {
			log.WithError(err).Error("Unmarshal failed")
			return nil, errors2.ErrInternal
		}
	}

	return result, nil
}
