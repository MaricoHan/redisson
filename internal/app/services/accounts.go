package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/account"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"time"
)

type IAccount interface {
	CreateAccount(account dto.CreateAccount) (*dto.AccountRes, error)
	GetAccounts(account dto.AccountsInfo) (*dto.AccountsRes, error)
}

type account struct {
	logger *log.Logger
}

func NewAccount(logger *log.Logger) *account {
	return &account{logger: logger}
}

func (a *account) CreateAccount(params dto.CreateAccount) (*dto.AccountRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "account"
	logFields["func"] = "CreateAccount"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	req := pb.AccountCreateRequest{
		ProjectId:  params.ProjectID,
		Count:      params.Count,
		PlatformId: int64(params.PlatFormID),
	}
	resp := &pb.AccountCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.AccountClientMap[mapKey]
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
	result := &dto.AccountRes{}
	result.Accounts = resp.Accounts
	return result, nil
}

func (a *account) GetAccounts(params dto.AccountsInfo) (*dto.AccountsRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "account"
	logFields["func"] = "GetAccounts"
	logFields["module"] = params.Module
	logFields["code"] = params.Code

	sort, ok := pb.Sorts_value[params.SortBy]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrSortBy)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrSortBy)
	}

	req := pb.AccountShowRequest{
		ProjectId: params.ProjectID,
		Offset:    params.Offset,
		Limit:     params.Limit,
		SortBy:    pb.Sorts(sort),
		Address:   params.Account,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	}

	resp := &pb.AccountShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.AccountClientMap[mapKey]
	if !ok {
		if !ok {
			log.WithFields(logFields).Error(errors2.ErrService)
			return nil, errors2.New(errors2.InternalError, errors2.ErrService)
		}
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
	result := &dto.AccountsRes{
		Accounts: []*dto.Account{},
	}
	result.Offset = resp.Offset
	result.Limit = resp.Limit
	result.TotalCount = resp.TotalCount
	var accounts []*dto.Account
	for _, result := range resp.Data {
		account := &dto.Account{
			Account: result.Address,
			Gas:     result.Gas,
			BizFee:  result.BizFee,
		}
		accounts = append(accounts, account)
	}

	if accounts != nil {
		result.Accounts = accounts
	}

	return result, nil
}
