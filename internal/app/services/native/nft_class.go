package native

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/class"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INFTClass interface {
	GetAllNFTClasses(ctx context.Context, params nft.NftClasses) (*nft.NftClassesRes, error) // 列表
	GetNFTClass(ctx context.Context, params nft.NftClasses) (*nft.NftClassRes, error)        // 详情
	CreateNFTClass(ctx context.Context, params nft.CreateNftClass) (*nft.TxRes, error)       // 创建
}

type nftClass struct {
	logger *log.Logger
}

func NewNFTClass(logger *log.Logger) *nftClass {
	return &nftClass{logger: logger}
}

func (n *nftClass) GetAllNFTClasses(ctx context.Context, params nft.NftClasses) (*nft.NftClassesRes, error) {
	logger := n.logger.WithField("params", params).WithField("func", "NFTClassList")

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
		Limit:      params.Limit,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		SortBy:     pb.SORTS(sort),
		Id:         params.Id,
		Name:       params.Name,
		Owner:      params.Owner,
		TxHash:     params.TxHash,
		CountTotal: params.CountTotal,
	}
	resp := &pb.ClassListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeNFTClassClientMap[mapKey]
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
	result := &nft.NftClassesRes{
		PageRes: dto.PageRes{
			PrevPageKey: resp.PrevPageKey,
			NextPageKey: resp.NextPageKey,
			TotalCount:  resp.TotalCount,
			Limit:       resp.Limit,
		},
		NftClasses: []*nft.NftClass{},
	}
	var nftClasses []*nft.NftClass
	for _, item := range resp.Data {
		nftClass := &nft.NftClass{
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

func (n *nftClass) GetNFTClass(ctx context.Context, params nft.NftClasses) (*nft.NftClassRes, error) {
	logger := n.logger.WithField("params", params).WithField("func", "GetNFTClass")

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
	grpcClient, ok := initialize.NativeNFTClassClientMap[mapKey]
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
	result := &nft.NftClassRes{}
	result.Id = resp.Detail.ClassId
	result.Timestamp = resp.Detail.Timestamp
	result.Name = resp.Detail.Name
	result.Uri = resp.Detail.Uri
	result.Owner = resp.Detail.Owner
	result.Symbol = resp.Detail.Symbol
	result.Data = resp.Detail.Metadata
	result.Description = resp.Detail.Description
	result.UriHash = resp.Detail.UriHash
	result.NftCount = resp.Detail.NftCount
	result.TxHash = resp.Detail.TxHash
	result.EditableByOwner = resp.Detail.UpdateRestricted
	return result, nil
}

func (n *nftClass) CreateNFTClass(ctx context.Context, params nft.CreateNftClass) (*nft.TxRes, error) {
	logger := n.logger.WithField("params", params).WithField("func", "CreateNFTClass")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.ClassCreateRequest{
		Name:             params.Name,
		Symbol:           params.Symbol,
		Description:      params.Description,
		Uri:              params.Uri,
		UriHash:          params.UriHash,
		Owner:            params.Owner,
		Metadata:         params.Data,
		ProjectId:        params.ProjectID,
		OperationId:      params.OperationId,
		ClassId:          params.ClassId,
		UpdateRestricted: params.EditableByOwner,
	}

	resp := &pb.ClassCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeNFTClassClientMap[mapKey]
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
	return &nft.TxRes{}, nil
}
