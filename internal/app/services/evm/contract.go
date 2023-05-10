package evm

import (
	"context"
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/evm"
	"time"

	log "github.com/sirupsen/logrus"

	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/contract"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IContract interface {
	CreateCall(ctx context.Context, params evm.CreateContractCall) (*evm.TxRes, error)
	ShowCall(ctx context.Context, params evm.ShowContractCall) (*evm.ShowContractCallRes, error)
}

type contract struct {
	logger *log.Entry
}

func NewContract(logger *log.Logger) *contract {
	return &contract{logger: logger.WithField("model", "contract")}
}

func (t *contract) CreateCall(ctx context.Context, params evm.CreateContractCall) (*evm.TxRes, error) {
	logger := t.logger.WithContext(ctx).WithField("params", params).WithField("func", "CreateCall")
	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.CreateCallRequest{
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
		From:        params.From,
		To:          params.To,
		Data:        params.Data,
		GasLimit:    params.GasLimit,
		Estimation:  params.Estimation,
	}
	resp := &pb.CreateCallResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.EvmContractClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.CreateCall(ctx, &req)
	if err != nil {
		logger.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &evm.TxRes{}, nil
}

func (t *contract) ShowCall(ctx context.Context, params evm.ShowContractCall) (*evm.ShowContractCallRes, error) {
	logger := t.logger.WithContext(ctx).WithField("params", params).WithField("func", "ShowCall")
	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.ShowCallRequest{
		From: params.From,
		To:   params.To,
		Data: params.Data,
	}
	resp := &pb.ShowCallResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.EvmContractClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.ShowCall(ctx, &req)
	if err != nil {
		logger.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	result := &evm.ShowContractCallRes{
		Result: resp.Result,
	}

	return result, nil
}
