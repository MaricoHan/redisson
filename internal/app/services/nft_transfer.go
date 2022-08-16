package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/class"
	pb2 "gitlab.bianjie.ai/avata/chains/api/pb/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"time"
)

type INFTTransfer interface {
	TransferNFTClass(params dto.TransferNftClassById) (*dto.TxRes, error) // 转让NFTClass
	TransferNFT(params dto.TransferNftByNftId) (*dto.TxRes, error)
}

type nftTransfer struct {
	logger *log.Logger
}

func NewNFTTransfer(logger *log.Logger) *nftTransfer {
	return &nftTransfer{logger: logger}
}

func (s *nftTransfer) TransferNFTClass(params dto.TransferNftClassById) (*dto.TxRes, error) {
	logger := s.logger.WithField("params",params).WithField("func","TransferNFTClass")

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := pb.ClassTransferRequest{
		ClassId:     params.ClassID,
		Owner:       params.Owner,
		Recipient:   params.Recipient,
		ProjectId:   params.ProjectID,
		Tag:         string(params.Tag),
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
	return &dto.TxRes{TaskId: resp.TaskId, OperationId: resp.OperationId}, nil

}

func (s *nftTransfer) TransferNFT(params dto.TransferNftByNftId) (*dto.TxRes, error) {
	logger := s.logger.WithField("params",params).WithField("func","TransferNFT")

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb2.NFTTransferRequest{
		ClassId:     params.ClassID,
		Owner:       params.Sender,
		NftId:       params.NftId,
		Recipient:   params.Recipient,
		ProjectId:   params.ProjectID,
		Tag:         string(params.Tag),
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
	return &dto.TxRes{TaskId: resp.TaskId, OperationId: resp.OperationId}, nil
}
