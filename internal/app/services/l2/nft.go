package l2

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/l2/nft"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/l2"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INFT interface {
	List(ctx context.Context, params dto.Nfts) (*dto.NftsRes, error)
	Create(ctx context.Context, params dto.CreateNfts) (*dto.TxRes, error)
	Show(ctx context.Context, params dto.NftByNftId) (*dto.NftRes, error)
	Update(ctx context.Context, params dto.EditNftByNftId) (*dto.TxRes, error)
	Delete(ctx context.Context, params dto.DeleteNftByNftId) (*dto.TxRes, error)
	Transfer(ctx context.Context, params dto.TransferNftByNftId) (*dto.TxRes, error)
}
type nft struct {
	logger *log.Logger
}

func NewNFT(logger *log.Logger) *nft {
	return &nft{logger: logger}
}

func (s *nft) List(ctx context.Context, params dto.Nfts) (*dto.NftsRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "NFTList")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	sort, ok := pb.SORTS_value[params.SortBy]
	if !ok {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.NFTListRequest{
		ProjectId:  params.ProjectID,
		Limit:      params.Limit,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		Id:         params.Id,
		ClassId:    params.ClassId,
		Owner:      params.Owner,
		TxHash:     params.TxHash,
		Status:     pb.STATUS(params.Status),
		SortBy:     pb.SORTS(sort),
		PageKey:    params.PageKey,
		CountTotal: params.CountTotal,
	}

	resp := &pb.NFTListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.List(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftsRes{
		Nfts: []*dto.NFT{},
	}
	result.Limit = resp.Limit
	result.PrevPageKey = resp.PrevPageKey
	result.NextPageKey = resp.NextPageKey
	result.TotalCount = resp.TotalCount

	var nfts []*dto.NFT
	for _, item := range resp.Data {
		nft := &dto.NFT{
			Id:          item.NftId,
			ClassId:     item.ClassId,
			Uri:         item.Uri,
			UriHash:     item.UriHash,
			Owner:       item.Owner,
			Status:      pb.STATUS_value[item.Status.String()],
			TxHash:      item.TxHash,
			Timestamp:   item.Timestamp,
			ClassName:   item.ClassName,
			ClassSymbol: item.ClassSymbol,
		}
		nfts = append(nfts, nft)
	}
	if nfts != nil {
		result.Nfts = nfts
	}

	return result, nil

}

func (s *nft) Create(ctx context.Context, params dto.CreateNfts) (*dto.TxRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "CreateNFT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.NFTCreateRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassId,
		Name:        params.Name,
		Uri:         params.Uri,
		UriHash:     params.UriHash,
		Data:        params.Data,
		Recipient:   params.Recipient,
		OperationId: params.OperationId,
	}
	resp := &pb.NFTCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Create(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.TxRes{}, nil

}

func (s *nft) Show(ctx context.Context, params dto.NftByNftId) (*dto.NftRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "ShowNFT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.NFTShowRequest{
		ProjectId: params.ProjectID,
		ClassId:   params.ClassId,
		NftId:     params.NftId,
	}
	resp := &pb.NFTShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftRes{
		Id:          resp.Detail.NftId,
		ClassId:     resp.Detail.ClassId,
		ClassName:   resp.Detail.ClassName,
		ClassSymbol: resp.Detail.ClassSymbol,
		Uri:         resp.Detail.Uri,
		UriHash:     resp.Detail.UriHash,
		Owner:       resp.Detail.Owner,
		Status:      pb.STATUS_value[resp.Detail.Status.String()],
		TxHash:      resp.Detail.TxHash,
		Timestamp:   resp.Detail.Timestamp,
	}

	return result, nil

}

func (s *nft) Update(ctx context.Context, params dto.EditNftByNftId) (*dto.TxRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "UpdateNFT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.NFTUpdateRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassId,
		Uri:         params.Uri,
		UriHash:     params.UriHash,
		NftId:       params.NftId,
		Owner:       params.Sender,
		OperationId: params.OperationId,
	}
	resp := &pb.NFTUpdateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Update(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("grpc response is nil")
	}
	return &dto.TxRes{}, nil
}

func (s *nft) Transfer(ctx context.Context, params dto.TransferNftByNftId) (*dto.TxRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "TransferNFT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.NFTTransferRequest{
		ClassId:     params.ClassID,
		Owner:       params.Sender,
		NftId:       params.NftId,
		Recipient:   params.Recipient,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}
	resp := &pb.NFTTransferResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClientMap[mapKey]
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

func (s *nft) Delete(ctx context.Context, params dto.DeleteNftByNftId) (*dto.TxRes, error) {
	logger := s.logger.WithContext(ctx).WithField("params", params).WithField("func", "DeleteNFT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.NFTDeleteRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassId,
		NftId:       params.NftId,
		Owner:       params.Sender,
		OperationId: params.OperationId,
	}
	resp := &pb.NFTDeleteResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Delete(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.TxRes{}, nil
}
