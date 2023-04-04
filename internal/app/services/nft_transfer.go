package services

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/class_v2"
	pb2 "gitlab.bianjie.ai/avata/chains/api/v2/pb/nft_v2"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INFTTransfer interface {
	TransferNFTClass(ctx context.Context, params dto.TransferNftClassById) (*dto.TxRes, error)
	TransferNFT(ctx context.Context, params dto.TransferNftByNftId) (*dto.TxRes, error)
}

type nftTransfer struct {
	logger *log.Logger
}

func NewNFTTransfer(logger *log.Logger) *nftTransfer {
	return &nftTransfer{logger: logger}
}

func (s *nftTransfer) TransferNFTClass(ctx context.Context, params dto.TransferNftClassById) (*dto.TxRes, error) {
	logger := s.logger.WithField("params", params).WithField("func", "TransferNFTClass")

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
	grpcClient, ok := initialize.ClassClientMap[mapKey]
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
	return &dto.TxRes{}, nil

}

func (s *nftTransfer) TransferNFT(ctx context.Context, params dto.TransferNftByNftId) (*dto.TxRes, error) {
	logger := s.logger.WithField("params", params).WithField("func", "TransferNFT")

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
	grpcClient, ok := initialize.NftClientMap[mapKey]
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
	return &dto.TxRes{}, nil
}
