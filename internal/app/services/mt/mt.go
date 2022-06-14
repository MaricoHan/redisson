package mt

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"time"

	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/mt"
)

type IMT interface {
	Issue(params *dto.IssueRequest) (*dto.IssueResponse, error)
	Mint(params *dto.MintRequest) (*dto.MintResponse, error)
	Show(params *dto.MTShowRequest) (*dto.MTShowResponse, error)
	List(params *dto.MTListRequest) (*dto.MTListResponse, error)
}
type MT struct {
	logger *log.Entry
}

func NewMT(logger *log.Logger) *MT {
	return &MT{
		logger: logger.WithField("service", "mt"),
	}
}

func (M MT) Issue(params *dto.IssueRequest) (*dto.IssueResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Issue"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	log := M.logger.WithFields(logFields)

	req := pb.MTIssueRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassID,
		Metadata:    params.Metadata,
		Recipients:  params.Recipients,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTIssueResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Issue(ctx, &req)
	if err != nil {
		log.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.IssueResponse{OperationID: params.OperationID}, nil
}

func (M MT) Mint(params *dto.MintRequest) (*dto.MintResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Issue"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	log := M.logger.WithFields(logFields)

	req := pb.MTMintRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassID,
		MtId:        params.MTID,
		Recipients:  params.Recipients,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTMintResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Mint(ctx, &req)
	if err != nil {
		log.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.MintResponse{OperationID: params.OperationID}, nil
}

func (M MT) Show(params *dto.MTShowRequest) (*dto.MTShowResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Show"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	req := pb.MTShowRequest{
		ProjectId: params.ProjectID,
		MtClassId: params.ClassID,
		MtId:      params.MTID,
	}
	resp := &pb.MTShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
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
	result := &dto.MTShowResponse{
		MtId:        resp.Data.MtId,
		MtClassId:   resp.Data.MtClassId,
		MtClassName: resp.Data.MtClassName,
		Data:        resp.Data.Data,
		OwnerCount:  resp.Data.OwnerCount,
		IssueData:   resp.Data.IssueData,
		MtCount:     resp.Data.MtCount,
		MintCount:   resp.Data.MintCount,
	}
	return result, nil
}

func (M MT) List(params *dto.MTListRequest) (*dto.MTListResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "List"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	sort, ok := pb.Sorts_value[params.SortBy]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.MTListRequest{
		ProjectId: params.ProjectID,
		Offset:    params.Offset,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		SortBy:    pb.Sorts(sort),
		MtId:      params.MtId,
		MtClassId: params.MtClassId,
		Issuer:    params.Issuer,
		TxHash:    params.TxHash,
	}

	resp := &pb.MTListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
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
	result := &dto.MTListResponse{
		Mts: []*dto.MT{},
	}
	result.Limit = resp.Limit
	result.Offset = resp.Offset
	result.TotalCount = resp.TotalCount
	var mts []*dto.MT
	for _, item := range resp.Data {
		mt := &dto.MT{
			MtId:        item.MtId,
			MtClassId:   item.MtClassId,
			MtClassName: item.MtClassName,
			Issuer:      item.Issuer,
			MtCount:     item.MtCount,
			OwnerCount:  item.OwnerCount,
			Timestamp:   item.Timestamp,
		}
		mts = append(mts, mt)
	}
	if mts != nil {
		result.Mts = mts
	}

	return result, nil
}
