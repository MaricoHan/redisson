package services

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	pb "gitlab.bianjie.ai/avata/chains/api/pb/v2/account"
	"gitlab.bianjie.ai/avata/chains/api/v2/pb/wallet"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IAccount interface {
	BatchCreateAccount(ctx context.Context, account dto.BatchCreateAccount) (*dto.BatchAccountRes, error)
	CreateAccount(ctx context.Context, account dto.CreateAccount) (*dto.AccountRes, error)
	GetAccounts(ctx context.Context, account dto.AccountsInfo) (*dto.AccountsRes, error)
	GetUserAccounts(ctx context.Context, account dto.AccountsInfo) (*dto.ShowUsersAccountsRes, error)
}

type account struct {
	logger *log.Logger
}

func NewAccount(logger *log.Logger) *account {
	return &account{logger: logger}
}

// BatchCreateAccount 批量创建链账户
func (a *account) BatchCreateAccount(ctx context.Context, params dto.BatchCreateAccount) (*dto.BatchAccountRes, error) {
	logger := a.logger.WithContext(ctx).WithField("params", params).WithField("func", "BatchCreateAccount")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	req := pb.AccountBatchCreateRequest{
		ProjectId:   params.ProjectID,
		Count:       params.Count,
		OperationId: params.OperationId,
	}
	resp := &pb.AccountBatchCreateResponse{}
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
	return result, nil
}

// CreateAccount 单个创建链账户
func (a *account) CreateAccount(ctx context.Context, params dto.CreateAccount) (*dto.AccountRes, error) {
	logger := a.logger.WithContext(ctx).WithField("params", params).WithField("func", "CreateAccount")

	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)

	if mapKey == constant.WalletServer {
		req := wallet.AccountCreateRequest{
			ProjectId:   params.ProjectID,
			Name:        params.Name,
			OperationId: params.OperationId,
			UserId:      params.UserId,
		}
		resp := &wallet.AccountCreateResponse{}
		var err error

		grpcClient, ok := initialize.WalletClientMap[mapKey]
		if !ok {
			logger.Error(errors2.ErrService)
			return nil, errors2.New(errors2.InternalError, errors2.ErrService)
		}
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
		defer cancel()
		resp, err = grpcClient.CreateAccount(ctx, &req)
		if err != nil {
			logger.WithError(err).Error("request err")
			return nil, err
		}
		if resp == nil {
			return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
		}
		result := &dto.AccountRes{}
		result.Account = resp.Account
		return result, nil
	}

	req := pb.AccountCreateRequest{
		ProjectId:   params.ProjectID,
		Name:        params.Name,
		OperationId: params.OperationId,
	}
	resp := &pb.AccountCreateResponse{}
	var err error

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
	return result, nil
}

// GetAccounts 查询链账户
func (a *account) GetAccounts(ctx context.Context, params dto.AccountsInfo) (*dto.AccountsRes, error) {
	logger := a.logger.WithContext(ctx).WithField("params", params).WithField("func", "GetAccounts")

	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)

	sort, ok := pb.SORT_value[params.SortBy]
	if !ok {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.AccountShowRequest{
		ProjectId:   params.ProjectID,
		PageKey:     params.PageKey,
		CountTotal:  params.CountTotal,
		Limit:       params.Limit,
		SortBy:      pb.SORT(sort),
		Account:     params.Account,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		OperationId: params.OperationId,
		Name:        params.Name,
	}

	resp := &pb.AccountShowResponse{}
	var err error

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
		PageRes: dto.PageRes{
			PrevPageKey: resp.PrevPageKey,
			NextPageKey: resp.NextPageKey,
			Limit:       resp.Limit,
			TotalCount:  resp.TotalCount,
		},
		Accounts: []*dto.Account{},
	}
	var accounts []*dto.Account
	for _, result := range resp.Data {
		account := &dto.Account{
			Account:     result.Address,
			Name:        result.Name,
			OperationId: result.OperationId,
		}
		accounts = append(accounts, account)
	}

	if accounts != nil {
		result.Accounts = accounts
	}

	return result, nil
}

// GetUserAccounts
//  @Description: 查询钱包用户
//  @receiver a
//  @param ctx
//  @param params
//  @return *dto.ShowUsersAccountsRes
//  @return error
func (a *account) GetUserAccounts(ctx context.Context, params dto.AccountsInfo) (*dto.ShowUsersAccountsRes, error) {
	logger := a.logger.WithContext(ctx).WithField("params", params).WithField("func", "GetUserAccounts")
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	sort, ok := wallet.SORT_value[params.SortBy]
	if !ok {
		logger.Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := wallet.AccountShowRequest{
		ProjectId:   params.ProjectID,
		PageKey:     params.PageKey,
		CountTotal:  params.CountTotal,
		Limit:       params.Limit,
		SortBy:      wallet.SORT(sort),
		Account:     params.Account,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		OperationId: params.OperationId,
		Name:        params.Name,
		UserId:      params.UserId,
	}

	resp := &wallet.AccountShowResponse{}
	var err error

	grpcClient, ok := initialize.WalletClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.ShowAccounts(ctx, &req)
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.ShowUsersAccountsRes{
		PageRes: dto.PageRes{
			PrevPageKey: resp.PrevPageKey,
			NextPageKey: resp.NextPageKey,
			Limit:       resp.Limit,
			TotalCount:  resp.TotalCount,
		},
		Accounts: []*dto.ShowUsersAccount{},
	}
	var accounts []*dto.ShowUsersAccount
	for _, v := range resp.Data {
		account := &dto.ShowUsersAccount{
			Account: dto.Account{
				Account:     v.Address,
				Name:        v.Name,
				OperationId: v.OperationId,
			},
			ReadOnly: uint32(v.ReadOnly),
		}
		accounts = append(accounts, account)
	}

	if accounts != nil {
		result.Accounts = accounts
	}

	return result, nil
}
