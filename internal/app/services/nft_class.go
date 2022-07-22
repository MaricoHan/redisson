package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/class"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"time"
)

type INFTClass interface {
	GetAllNFTClasses(params dto.NftClasses) (*dto.NftClassesRes, error) // 列表
	GetNFTClass(params dto.NftClasses) (*dto.NftClassRes, error)        // 详情
	CreateNFTClass(params dto.CreateNftClass) (*dto.TxRes, error)       // 创建
}

type nftClass struct {
	logger *log.Logger
}

func NewNFTClass(logger *log.Logger) *nftClass {
	return &nftClass{logger: logger}
}

func (n *nftClass) GetAllNFTClasses(params dto.NftClasses) (*dto.NftClassesRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nftclass"
	logFields["func"] = "GetAllNFTClasses"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	sort, ok := pb.SORTS_value[params.SortBy]
	if !ok {
		log.WithFields(logFields).Error("sort_by is illegal")
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.ClassListRequest{
		ProjectId: params.ProjectID,
		Offset:    params.Offset,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		SortBy:    pb.SORTS(sort),
		Id:        params.Id,
		Name:      params.Name,
		Owner:     params.Owner,
		TxHash:    params.TxHash,
		Status:    pb.STATUS_Active,
	}
	resp := &pb.ClassListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.ClassClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.List(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftClassesRes{
		PageRes: dto.PageRes{
			Offset: resp.Offset,
			Limit:  resp.Limit,
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
			NftCount:  item.NftCount,
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

func (n *nftClass) GetNFTClass(params dto.NftClasses) (*dto.NftClassRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nftclass"
	logFields["func"] = "GetNFTClass"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.ClassShowRequest{
		ProjectId: params.ProjectID,
		Id:        params.Id,
		Status:    pb.STATUS_Active, //todo
	}
	resp := &pb.ClassShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.ClassClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftClassRes{}
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
	return result, nil
}

func (n *nftClass) CreateNFTClass(params dto.CreateNftClass) (*dto.TxRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nftclass"
	logFields["func"] = "CreateNFTClass"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.ClassCreateRequest{
		Name:        params.Name,
		Symbol:      params.Symbol,
		Description: params.Description,
		Uri:         params.Uri,
		UriHash:     params.UriHash,
		Owner:       params.Owner,
		Data:        params.Data,
		ProjectId:   params.ProjectID,
		Tag:         string(params.Tag),
		OperationId: params.OperationId,
		ClassId:     params.ClassId,
	}

	resp := &pb.ClassCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.ClassClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Create(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.TxRes{TaskId: resp.TaskId, OperationId: resp.OperationId}, nil
}
