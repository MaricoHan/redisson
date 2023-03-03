package services

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/account"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
)

type IAccount interface {
	BatchCreateAccount(ctx context.Context, account dto.BatchCreateAccount) (*dto.BatchAccountRes, error)
	CreateAccount(ctx context.Context, account dto.CreateAccount) (*dto.AccountRes, error)
	GetAccounts(ctx context.Context, account dto.AccountsInfo) (*dto.AccountsRes, error)
}

type account struct {
	logger *log.Logger
}

func NewAccount(logger *log.Logger) *account {
	return &account{logger: logger}
}

// BatchCreateAccount 批量创建链账户
func (a *account) BatchCreateAccount(ctx context.Context, params dto.BatchCreateAccount) (*dto.BatchAccountRes, error) {
	logger := a.logger.WithField("params", params).WithField("func", "BatchCreateAccount")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.AccountCreateRequest{
		ProjectId:   params.ProjectID,
		Count:       params.Count,
		PlatformId:  int64(params.PlatFormID),
		OperationId: params.OperationId,
	}
	resp := &pb.AccountCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.AccountClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.BatchCreate(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.BatchAccountRes{}
	result.Accounts = resp.Accounts
	result.OperationId = resp.OperationId
	return result, nil
}

// CreateAccount 单个创建链账户
func (a *account) CreateAccount(ctx context.Context, params dto.CreateAccount) (*dto.AccountRes, error) {
	logger := a.logger.WithField("params", params).WithField("func", "CreateAccount")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.AccountSeparateCreateRequest{
		ProjectId:   params.ProjectID,
		Name:        params.Name,
		PlatformId:  int64(params.PlatFormID),
		OperationId: params.OperationId,
	}
	resp := &pb.AccountSeparateCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.AccountClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.Create(ctx, &req)
	if err != nil {
		logger.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.AccountRes{}
	result.Account = resp.Account
	result.Name = resp.Name
	result.OperationId = resp.OperationId
	return result, nil
}

func (a *account) GetAccounts(ctx context.Context, params dto.AccountsInfo) (*dto.AccountsRes, error) {
	logger := a.logger.WithField("params", params).WithField("func", "GetAccounts")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	sort, ok := pb.SORT_value[params.SortBy]
	if !ok {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.AccountShowRequest{
		ProjectId:   params.ProjectID,
		NextKey:     params.NextKey,
		CountTotal:  params.CountTotal,
		Limit:       params.Limit,
		SortBy:      pb.SORT(sort),
		Address:     params.Account,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		OperationId: params.OperationId,
		Name:        params.Name,
	}

	resp := &pb.AccountShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.AccountClientMap[mapKey]
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
	result := &dto.AccountsRes{
		Accounts: []*dto.Account{},
	}
	result.Offset = resp.Offset
	result.Limit = resp.Limit
	result.TotalCount = resp.TotalCount
	var accounts []*dto.Account
	for _, result := range resp.Data {
		account := &dto.Account{
			Account:     result.Address,
			Gas:         result.Gas,
			BizFee:      result.BizFee,
			Name:        result.Name,
			OperationId: result.OperationId,
			Status:      uint64(result.Status),
		}
		accounts = append(accounts, account)
	}

	if accounts != nil {
		result.Accounts = accounts
	}

	return result, nil
}
