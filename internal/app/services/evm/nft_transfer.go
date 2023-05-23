package evm

import (
	"context"
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/evm"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/class"
	pb2 "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INFTTransfer interface {
	TransferNFTClass(ctx context.Context, params evm.TransferNftClassById) (*evm.TxRes, error)
	TransferNFT(ctx context.Context, params evm.TransferNftByNftId) (*evm.TxRes, error)
}

type nftTransfer struct {
	logger *log.Logger
}

func NewNFTTransfer(logger *log.Logger) *nftTransfer {
	return &nftTransfer{logger: logger}
}

func (s *nftTransfer) TransferNFTClass(ctx context.Context, params evm.TransferNftClassById) (*evm.TxRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "TransferNFTClass")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := pb.ClassTransferRequest{
		ClassId:     params.ClassID,
		Owner:       params.Owner,
		Recipient:   params.Recipient,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}
	resp := &pb.ClassTransferResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.EvmClassClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Transfer(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &evm.TxRes{}, nil

}

func (s *nftTransfer) TransferNFT(ctx context.Context, params evm.TransferNftByNftId) (*evm.TxRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "TransferNFT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb2.NFTTransferRequest{
		ClassId:     params.ClassID,
		Owner:       params.Sender,
		NftId:       params.NftId,
		Recipient:   params.Recipient,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}
	resp := &pb2.NFTTransferResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.EvmNftClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Transfer(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &evm.TxRes{}, nil
}
