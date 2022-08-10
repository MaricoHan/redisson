package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/nft"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
)

type INFT interface {
	List(params dto.Nfts) (*dto.NftsRes, error)
	Create(params dto.CreateNfts) (*dto.TxRes, error)
	BatchCreate(params dto.BatchCreateNfts) (*dto.BatchTxRes, error)
	Show(params dto.NftByNftId) (*dto.NftReq, error)
	Update(params dto.EditNftByNftId) (*dto.TxRes, error)
	Delete(params dto.DeleteNftByNftId) (*dto.TxRes, error)
	BatchTransfer(params *dto.BatchTransferRequest) (*dto.BatchTxRes, error)
	BatchEdit(params *dto.BatchEditRequest) (*dto.BatchTxRes, error)
	BatchDelete(params *dto.BatchDeleteRequest) (*dto.BatchTxRes, error)
}
type nft struct {
	logger *log.Logger
}

func NewNFT(logger *log.Logger) *nft {
	return &nft{logger: logger}
}

func (s *nft) List(params dto.Nfts) (*dto.NftsRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "List"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	sort, ok := pb.SORTS_value[params.SortBy]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.NFTListRequest{
		ProjectId: params.ProjectID,
		Offset:    params.Offset,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		Id:        params.Id,
		ClassId:   params.ClassId,
		Owner:     params.Owner,
		TxHash:    params.TxHash,
		Status:    pb.STATUS(pb.STATUS_value[params.Status]),
		SortBy:    pb.SORTS(sort),
		Name:      params.Name,
	}

	resp := &pb.NFTListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.List(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftsRes{
		Nfts: []*dto.NFT{},
	}
	result.Limit = resp.Limit
	result.Offset = resp.Offset
	result.TotalCount = resp.TotalCount
	var nfts []*dto.NFT
	for _, item := range resp.Data {
		nft := &dto.NFT{
			Id:          item.NftId,
			Name:        item.Name,
			ClassId:     item.ClassId,
			Uri:         item.Uri,
			Owner:       item.Owner,
			Status:      item.Status.String(),
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

func (s *nft) Create(params dto.CreateNfts) (*dto.TxRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "Create"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.NFTCreateRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassId,
		Name:        params.Name,
		Uri:         params.Uri,
		UriHash:     params.UriHash,
		Data:        params.Data,
		Recipient:   params.Recipient,
		Tag:         string(params.Tag),
		OperationId: params.OperationId,
	}
	resp := &pb.NFTCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
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

func (s *nft) BatchCreate(params dto.BatchCreateNfts) (*dto.BatchTxRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "BatchCreate"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.NFTBatchCreateRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassId,
		Name:        params.Name,
		Uri:         params.Uri,
		UriHash:     params.UriHash,
		Data:        params.Data,
		Recipients:  params.Recipients,
		Tag:         string(params.Tag),
		OperationId: params.OperationId,
	}
	resp := &pb.NFTBatchCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchCreate(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.BatchTxRes{OperationId: resp.OperationId}, nil
}

func (s *nft) Show(params dto.NftByNftId) (*dto.NftReq, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "Update"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	req := pb.NFTShowRequest{
		ProjectId: params.ProjectID,
		ClassId:   params.ClassId,
		NftId:     params.NftId,
	}
	resp := &pb.NFTShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.NftReq{
		Id:          resp.Detail.NftId,
		Name:        resp.Detail.Name,
		ClassId:     resp.Detail.ClassId,
		ClassName:   resp.Detail.ClassName,
		ClassSymbol: resp.Detail.ClassSymbol,
		Uri:         resp.Detail.Uri,
		UriHash:     resp.Detail.UriHash,
		Data:        resp.Detail.Metadata,
		Owner:       resp.Detail.Owner,
		Status:      resp.Detail.Status.String(),
		TxHash:      resp.Detail.TxHash,
		Timestamp:   resp.Detail.Timestamp,
	}

	return result, nil

}

func (s *nft) Update(params dto.EditNftByNftId) (*dto.TxRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "Update"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.NFTUpdateRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassId,
		Name:        params.Name,
		Uri:         params.Uri,
		NftId:       params.NftId,
		Data:        params.Data,
		Owner:       params.Sender,
		Tag:         string(params.Tag),
		OperationId: params.OperationId,
	}
	resp := &pb.NFTUpdateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Update(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("grpc response is nil")
	}
	return &dto.TxRes{TaskId: resp.TaskId, OperationId: resp.OperationId}, nil
}

func (s *nft) Delete(params dto.DeleteNftByNftId) (*dto.TxRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "Delete"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.NFTDeleteRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassId,
		NftId:       params.NftId,
		Owner:       params.Sender,
		Tag:         string(params.Tag),
		OperationId: params.OperationId,
	}
	resp := &pb.NFTDeleteResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Delete(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.TxRes{TaskId: resp.TaskId, OperationId: resp.OperationId}, nil
}
func (s *nft) BatchTransfer(params *dto.BatchTransferRequest) (*dto.BatchTxRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "BatchTransfer"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.NFTBatchTransferRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Sender,
		Data:        params.Data,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}
	resp := new(pb.NFTBatchTransferResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchTransfer(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.BatchTxRes{OperationId: resp.OperationId}, nil
}
func (s *nft) BatchEdit(params *dto.BatchEditRequest) (*dto.BatchTxRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "BatchEdit"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.NFTBatchEditRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Sender,
		Nfts:        params.Nfts,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}
	resp := new(pb.NFTBatchEditResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchEdit(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.BatchTxRes{OperationId: resp.OperationId}, nil
}

func (s *nft) BatchDelete(params *dto.BatchDeleteRequest) (*dto.BatchTxRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "nft"
	logFields["func"] = "BatchDelete"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.NFTBatchDeleteRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Sender,
		Nfts:        params.Nfts,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}
	resp := new(pb.NFTBatchDeleteResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NftClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchDelete(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.BatchTxRes{OperationId: resp.OperationId}, nil
}
