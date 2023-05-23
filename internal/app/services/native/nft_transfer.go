package native

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/class"
	pb2 "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INFTTransfer interface {
	TransferNFTClass(ctx context.Context, params nft.TransferNftClassById) (*nft.TxRes, error) // 转让NFTClass
	TransferNFT(ctx context.Context, params nft.TransferNftByNftId) (*nft.TxRes, error)
}

type nftTransfer struct {
	logger *log.Logger
}

func NewNFTTransfer(logger *log.Logger) *nftTransfer {
	return &nftTransfer{logger: logger}
}

func (s *nftTransfer) TransferNFTClass(ctx context.Context, params nft.TransferNftClassById) (*nft.TxRes, error) {
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
	grpcClient, ok := initialize.NativeNFTClassClientMap[mapKey]
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
	return &nft.TxRes{}, nil

}

func (s *nftTransfer) TransferNFT(ctx context.Context, params nft.TransferNftByNftId) (*nft.TxRes, error) {
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
	grpcClient, ok := initialize.NativeNFTClientMap[mapKey]
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
	return &nft.TxRes{}, nil
}
