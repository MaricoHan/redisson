package native

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/mt_class"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IMTClass interface {
	Show(ctx context.Context, params *dto.MTClassShowRequest) (*dto.MTClassShowResponse, error)
	List(ctx context.Context, params *dto.MTClassListRequest) (*dto.MTClassListResponse, error)
	CreateMTClass(ctx context.Context, params dto.CreateMTClass) (*dto.BatchTxRes, error)     // 创建
	TransferMTClass(ctx context.Context, params dto.TransferMTClass) (*dto.BatchTxRes, error) // 转让
}

type mtClass struct {
	logger *log.Logger
}

func NewMTClass(logger *log.Logger) *mtClass {
	return &mtClass{logger: logger}
}

func (m *mtClass) CreateMTClass(ctx context.Context, params dto.CreateMTClass) (*dto.BatchTxRes, error) {
	logger := m.logger.WithField("params", params).WithField("func", "CreateMTClass")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.MTClassCreateRequest{
		Name:        params.Name,
		Owner:       params.Owner,
		Data:        params.Data,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}

	resp := &pb.MTClassCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClassClientMap[mapKey]
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
	return &dto.BatchTxRes{}, nil
}

func (m *mtClass) Show(ctx context.Context, params *dto.MTClassShowRequest) (*dto.MTClassShowResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "ShowMTClass")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTClassShowRequest{
		ProjectId: params.ProjectID,
		ClassId:   params.ClassID,
	}
	resp := &pb.MTClassShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClassClientMap[mapKey]
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
	result := &dto.MTClassShowResponse{
		// Id:          resp.Detail.Id,
		ClassId:   resp.Detail.ClassId,
		ClassName: resp.Detail.ClassName,
		Owner:     resp.Detail.Owner,
		Data:      resp.Detail.Data,
		// Status:      resp.Detail.Status,
		// LockedBy:    resp.Detail.LockedBy,
		TxHash:    resp.Detail.TxHash,
		Timestamp: resp.Detail.Timestamp,
		MtCount:   resp.Detail.MtCount,
		// Extra1:      resp.Detail.Extra1,
		// Extra2:      resp.Detail.Extra2,
		// CreatedAt:   resp.Detail.CreatedAt,
		// UpdatedAt:   resp.Detail.UpdatedAt,
	}
	return result, nil
}

func (m *mtClass) List(ctx context.Context, params *dto.MTClassListRequest) (*dto.MTClassListResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "MTClassList")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	sort, ok := pb.SORTS_value[params.SortBy]
	if !ok {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.MTClassListRequest{
		ProjectId: params.ProjectID,
		PageKey:   params.PageKey,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		SortBy:    pb.SORTS(sort),
		ClassId:   params.ClassId,
		ClassName: params.ClassName,
		Owner:     params.Owner,
		TxHash:    params.TxHash,
		//Status:    pb.STATUS(pb.STATUS_value[params.Status]),
	}

	resp := &pb.MTClassListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClassClientMap[mapKey]
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
	result := &dto.MTClassListResponse{
		MtClasses: []*dto.MTClass{},
	}
	result.Limit = resp.Limit
	result.PrevPageKey = resp.PrevPageKey
	result.NextPageKey = resp.NextPageKey
	result.TotalCount = resp.TotalCount
	var mtClasses []*dto.MTClass
	for _, item := range resp.Data {
		mtClass := &dto.MTClass{
			ClassId:   item.ClassId,
			ClassName: item.ClassName,
			Owner:     item.Owner,
			MtCount:   item.MtCount,
			TxHash:    item.TxHash,
			Timestamp: item.Timestamp,
		}
		mtClasses = append(mtClasses, mtClass)
	}
	if mtClasses != nil {
		result.MtClasses = mtClasses
	}

	return result, nil
}

func (m *mtClass) TransferMTClass(ctx context.Context, params dto.TransferMTClass) (*dto.BatchTxRes, error) {
	logger := m.logger.WithField("params", params).WithField("func", "TransferMTClass")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := pb.MTClassTransferRequest{
		ClassId:     params.ClassID,
		Owner:       params.Owner,
		Recipient:   params.Recipient,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}
	resp := &pb.MTClassTransferResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClassClientMap[mapKey]
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
	return &dto.BatchTxRes{}, nil
}
