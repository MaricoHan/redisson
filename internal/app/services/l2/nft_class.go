package l2

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/l2/nft"
	dto_1 "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/l2"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INFTClass interface {
	List(ctx context.Context, params dto.NftClasses) (*dto.NftClassesRes, error)       // 列表
	Show(ctx context.Context, params dto.NftClasses) (*dto.NftClassRes, error)         // 详情
	Create(ctx context.Context, params dto.CreateNftClass) (*dto.TxRes, error)         // 创建
	Transfer(ctx context.Context, params dto.TransferNftClassById) (*dto.TxRes, error) // 转让
}

type nftClass struct {
	logger *log.Logger
}

func NewNFTClass(logger *log.Logger) *nftClass {
	return &nftClass{logger: logger}
}

func (n *nftClass) List(ctx context.Context, params dto.NftClasses) (*dto.NftClassesRes, error) {
	logger := n.logger.WithContext(ctx).WithField("params", params).WithField("func", "L2NFTClassList")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	sort, ok := pb.SORTS_value[params.SortBy]
	if !ok {
		logger.Error("sort_by is illegal")
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.ClassListRequest{
		ProjectId:  params.ProjectID,
		PageKey:    params.PageKey,
		CountTotal: params.CountTotal,
		Limit:      params.Limit,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		SortBy:     pb.SORTS(sort),
		Id:         params.Id,
		Name:       params.Name,
		Owner:      params.Owner,
		TxHash:     params.TxHash,
	}
	resp := &pb.ClassListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClassClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.List(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftClassesRes{
		PageRes: dto_1.PageRes{
			PrevPageKey: resp.PrevPageKey,
			NextPageKey: resp.NextPageKey,
			Limit:       resp.Limit,
			TotalCount:  resp.TotalCount,
		},
		NftClasses: []*dto.NftClass{},
	}
	var nftClasses []*dto.NftClass
	for _, item := range resp.Data {
		nftClass := &dto.NftClass{
			Id:        item.ClassId,
			Name:      item.Name,
			Symbol:    item.Symbol,
			Uri:       item.Uri,
			Owner:     item.Owner,
			TxHash:    item.TxHash,
			Timestamp: item.Timestamp,
		}
		nftClasses = append(nftClasses, nftClass)
	}
	result.TotalCount = resp.TotalCount

	if nftClasses != nil {
		result.NftClasses = nftClasses
	}
	return result, nil
}

func (n *nftClass) Show(ctx context.Context, params dto.NftClasses) (*dto.NftClassRes, error) {
	logger := n.logger.WithContext(ctx).WithField("params", params).WithField("func", "GetL2NFTClass")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.ClassShowRequest{
		ProjectId: params.ProjectID,
		Id:        params.Id,
	}
	resp := &pb.ClassShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClassClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftClassRes{}
	result.Id = resp.Detail.ClassId
	result.Name = resp.Detail.Name
	result.Symbol = resp.Detail.Symbol
	result.Description = resp.Detail.Description
	result.Uri = resp.Detail.Uri
	result.UriHash = resp.Detail.Uri
	result.Data = resp.Detail.Data
	result.Owner = resp.Detail.Owner
	result.NftCount = resp.Detail.NftCount
	result.TxHash = resp.Detail.TxHash
	result.Timestamp = resp.Detail.Timestamp
	return result, nil
}

func (n *nftClass) Create(ctx context.Context, params dto.CreateNftClass) (*dto.TxRes, error) {
	logger := n.logger.WithContext(ctx).WithField("params", params).WithField("func", "CreateL2NFTClass")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.ClassCreateRequest{
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
		Name:        params.Name,
		ClassId:     params.ClassId,
		Symbol:      params.Symbol,
		Description: params.Description,
		Uri:         params.Uri,
		UriHash:     params.UriHash,
		Data:        params.Data,
		Owner:       params.Owner,
	}

	resp := &pb.ClassCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.L2NftClassClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
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

func (s *nftClass) Transfer(ctx context.Context, params dto.TransferNftClassById) (*dto.TxRes, error) {
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
	grpcClient, ok := initialize.L2NftClassClientMap[mapKey]
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