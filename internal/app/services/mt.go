package services

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IMT interface {
	Issue(params *dto.IssueRequest) (*dto.IssueResponse, error)
	Mint(params *dto.MintRequest) (*dto.MintResponse, error)
	Edit(params *dto.EditRequest) (*dto.EditResponse, error)
	Burn(params *dto.BurnRequest) (*dto.BurnResponse, error)
	Transfer(params *dto.MTTransferRequest) (*dto.MTTransferResponse, error)

	BatchTransfer(params *dto.MTBatchTransferRequest) (*dto.MTBatchTransferResponse, error)
	Show(params *dto.MTShowRequest) (*dto.MTShowResponse, error)
	List(params *dto.MTListRequest) (*dto.MTListResponse, error)
	Balances(params *dto.MTBalancesRequest) (*dto.MTBalancesResponse, error)
}
type MT struct {
	logger *log.Entry
}

func NewMT(logger *log.Logger) *MT {
	return &MT{
		logger: logger.WithField("service", "mt"),
	}
}

func (m MT) Issue(params *dto.IssueRequest) (*dto.IssueResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Issue"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	log := m.logger.WithFields(logFields)

	req := pb.MTIssueRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassID,
		Metadata:    params.Metadata,
		Amount:      params.Amount,
		Recipient:   params.Recipient,
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

func (m MT) Mint(params *dto.MintRequest) (*dto.MintResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Issue"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	log := m.logger.WithFields(logFields)

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

func (m MT) Show(params *dto.MTShowRequest) (*dto.MTShowResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Show"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	req := pb.MTShowRequest{
		ProjectId: params.ProjectID,
		ClassId:   params.ClassID,
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
		MtId:       resp.Data.MtId,
		ClassId:    resp.Data.ClassId,
		ClassName:  resp.Data.ClassName,
		Data:       resp.Data.Data,
		OwnerCount: resp.Data.OwnerCount,
		IssueData:  resp.Data.IssueData,
		MtCount:    resp.Data.MtCount,
		MintTimes:  resp.Data.MintCount,
	}
	return result, nil
}
func (m MT) Edit(params *dto.EditRequest) (*dto.EditResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Edit"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	log := m.logger.WithFields(logFields)

	req := pb.MTEditRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		Mts:         params.Mts,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTEditResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Edit(ctx, &req)
	if err != nil {
		log.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.EditResponse{OperationID: params.OperationID}, nil
}

func (m MT) Burn(params *dto.BurnRequest) (*dto.BurnResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Burn"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	log := m.logger.WithFields(logFields)

	req := pb.MTDeleteRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		Mts:         params.Mts,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTDeleteResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Delete(ctx, &req)
	if err != nil {
		log.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.BurnResponse{OperationID: params.OperationID}, nil
}
func (m MT) Transfer(params *dto.MTTransferRequest) (*dto.MTTransferResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "Transfer"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	log := m.logger.WithFields(logFields)

	req := pb.MTTransferRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		ClassId:     params.ClassId,
		MtId:        params.MtId,
		Amount:      params.Amount,
		Recipient:   params.Recipient,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTTransferResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Transfer(ctx, &req)
	if err != nil {
		log.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.MTTransferResponse{OperationID: params.OperationID}, nil
}

func (m MT) BatchTransfer(params *dto.MTBatchTransferRequest) (*dto.MTBatchTransferResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "BatchTransfer"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	log := m.logger.WithFields(logFields)

	req := pb.MTBatchTransferRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		Mts:         params.Mts,
		Tag:         params.Tag,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTBatchTransferResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchTransfer(ctx, &req)
	if err != nil {
		log.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.MTBatchTransferResponse{OperationID: params.OperationID}, nil
}
func (m MT) List(params *dto.MTListRequest) (*dto.MTListResponse, error) {
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
		ClassId:   params.ClassId,
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
			MtId:       item.MtId,
			ClassId:    item.ClassId,
			ClassName:  item.ClassName,
			Issuer:     item.Issuer,
			TxHash:     item.TxHash,
			OwnerCount: item.OwnerCount,
			Timestamp:  item.Timestamp,
		}
		mts = append(mts, mt)
	}
	if mts != nil {
		result.Mts = mts
	}

	return result, nil
}

func (M MT) Balances(params *dto.MTBalancesRequest) (*dto.MTBalancesResponse, error) {
	logFields := log.Fields{}
	logFields["model"] = "mt"
	logFields["func"] = "List"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.MTBalancesRequest{
		ProjectId: params.ProjectID,
		Offset:    params.Offset,
		Limit:     params.Limit,
		ClassId:   params.ClassId,
		Account:   params.Account,
		MtId:      params.MtId,
	}

	resp := &pb.MTBalancesResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.MTClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Balances(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.MTBalancesResponse{
		Mts: []*dto.MTBalances{},
	}
	result.Limit = resp.Limit
	result.Offset = resp.Offset
	result.TotalCount = resp.TotalCount
	var mts []*dto.MTBalances
	for _, item := range resp.Mts {
		mt := &dto.MTBalances{
			MtId:   item.MtId,
			Amount: item.Amount,
		}
		mts = append(mts, mt)
	}
	if mts != nil {
		result.Mts = mts
	}

	return result, nil
}
