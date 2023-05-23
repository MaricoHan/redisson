package native

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/mt"
	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IMT interface {
	Issue(ctx context.Context, params *dto.IssueRequest) (*dto.IssueResponse, error)
	Mint(ctx context.Context, params *dto.MintRequest) (*dto.MintResponse, error)
	BatchMint(ctx context.Context, params *dto.BatchMintRequest) (*dto.BatchMintResponse, error)

	Edit(ctx context.Context, params *dto.EditRequest) (*dto.EditResponse, error)
	Burn(ctx context.Context, params *dto.BurnRequest) (*dto.BurnResponse, error)
	Transfer(ctx context.Context, params *dto.MTTransferRequest) (*dto.MTTransferResponse, error)

	BatchTransfer(ctx context.Context, params *dto.MTBatchTransferRequest) (*dto.MTBatchTransferResponse, error)
	BatchBurn(ctx context.Context, params *dto.BatchBurnRequest) (*dto.BatchBurnResponse, error)
	Show(ctx context.Context, params *dto.MTShowRequest) (*dto.MTShowResponse, error)
	List(ctx context.Context, params *dto.MTListRequest) (*dto.MTListResponse, error)
	Balances(ctx context.Context, params *dto.MTBalancesRequest) (*dto.MTBalancesResponse, error)
}
type MT struct {
	logger *log.Entry
}

func NewMT(logger *log.Logger) *MT {
	return &MT{
		logger: logger.WithField("service", "mt"),
	}
}

func (m MT) Issue(ctx context.Context, params *dto.IssueRequest) (*dto.IssueResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "IssueMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTIssueRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassID,
		Metadata:    params.Metadata,
		Amount:      params.Amount,
		Recipient:   params.Recipient,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTIssueResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Issue(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.IssueResponse{}, nil
}

func (m MT) Mint(ctx context.Context, params *dto.MintRequest) (*dto.MintResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "MintMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTMintRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassID,
		MtId:        params.MTID,
		Amount:      params.Amount,
		Recipient:   params.Recipient,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTMintResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Mint(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.MintResponse{}, nil
}

func (m MT) BatchMint(ctx context.Context, params *dto.BatchMintRequest) (*dto.BatchMintResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "MintMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTBatchMintRequest{
		ProjectId:   params.ProjectID,
		ClassId:     params.ClassID,
		MtId:        params.MTID,
		Recipients:  params.Recipients,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTBatchMintResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchMint(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.BatchMintResponse{}, nil
}

func (m MT) Show(ctx context.Context, params *dto.MTShowRequest) (*dto.MTShowResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "ShowMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTShowRequest{
		ProjectId: params.ProjectID,
		ClassId:   params.ClassID,
		MtId:      params.MTID,
	}
	resp := &pb.MTShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
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

func (m MT) Edit(ctx context.Context, params *dto.EditRequest) (*dto.EditResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "EditMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTEditRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		Data:        params.Data,
		ClassId:     params.ClassId,
		MtId:        params.MTID,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTEditResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Edit(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.EditResponse{}, nil
}

func (m MT) BatchBurn(ctx context.Context, params *dto.BatchBurnRequest) (*dto.BatchBurnResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "BurnMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTBatchDeleteRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		Mts:         params.Mts,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTBatchDeleteResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchDelete(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.BatchBurnResponse{}, nil
}

func (m MT) Burn(ctx context.Context, params *dto.BurnRequest) (*dto.BurnResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "BurnMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTDeleteRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		ClassId:     params.ClassID,
		MtId:        params.MtID,
		Amount:      params.Amount,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTDeleteResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
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

	return &dto.BurnResponse{}, nil
}

func (m MT) Transfer(ctx context.Context, params *dto.MTTransferRequest) (*dto.MTTransferResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "TransferMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTTransferRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		ClassId:     params.ClassId,
		MtId:        params.MtId,
		Amount:      params.Amount,
		Recipient:   params.Recipient,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTTransferResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Transfer(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.MTTransferResponse{}, nil
}

func (m MT) BatchTransfer(ctx context.Context, params *dto.MTBatchTransferRequest) (*dto.MTBatchTransferResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "BatchTransferMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTBatchTransferRequest{
		ProjectId:   params.ProjectID,
		Owner:       params.Owner,
		Mts:         params.Mts,
		OperationId: params.OperationID,
	}

	resp := new(pb.MTBatchTransferResponse)

	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchTransfer(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.MTBatchTransferResponse{}, nil
}

func (m MT) List(ctx context.Context, params *dto.MTListRequest) (*dto.MTListResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "ListMT")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	sort, ok := pb.SORTS_value[params.SortBy]
	if !ok {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.MTListRequest{
		ProjectId: params.ProjectID,
		PageKey:   params.PageKey,
		Limit:     params.Limit,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		SortBy:    pb.SORTS(sort),
		MtId:      params.MtId,
		ClassId:   params.ClassId,
		Issuer:    params.Issuer,
		TxHash:    params.TxHash,
	}

	resp := &pb.MTListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
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
	result := &dto.MTListResponse{
		Mts: []*dto.MT{},
	}
	result.Limit = resp.Limit
	result.NextPageKey = resp.NextPageKey
	result.PrevPageKey = resp.PrevPageKey
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

func (m MT) Balances(ctx context.Context, params *dto.MTBalancesRequest) (*dto.MTBalancesResponse, error) {
	logger := m.logger.WithField("params", params).WithField("func", "BalancesList")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.MTBalancesRequest{
		ProjectId: params.ProjectID,
		PageKey:   params.PageKey,
		Limit:     params.Limit,
		ClassId:   params.ClassId,
		Account:   params.Account,
		MtId:      params.MtId,
	}

	resp := &pb.MTBalancesResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NativeMTClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Balances(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.MTBalancesResponse{
		Mts: []*dto.MTBalances{},
	}
	result.Limit = resp.Limit
	result.PrevPageKey = resp.PrevPageKey
	result.NextPageKey = resp.NextPageKey
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
