package services

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/wallet"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IUser interface {
	CreateUsers(ctx context.Context, params dto.CreateUsers) (*dto.CreateUsersRes, error)
	UpdateUsers(ctx context.Context, params dto.UpdateUsers) (*dto.TxRes, error)
	ShowUsers(ctx context.Context, params dto.ShowUsers) (*dto.CreateUsersRes, error)
}

type user struct {
	logger *log.Logger
}

func NewUser(logger *log.Logger) *user {
	return &user{logger: logger}
}

func (u user) CreateUsers(ctx context.Context, params dto.CreateUsers) (*dto.CreateUsersRes, error) {
	log := u.logger.WithContext(ctx).WithFields(
		map[string]interface{}{
			"function": "CreateUsers",
			"params":   params,
		})

	req := wallet.CreateUsersRequest{
		ProjectId:  params.ProjectID,
		UserType:   wallet.USER_TYPE(params.Usertype),
		Individual: params.Individual,
		Enterprise: params.Enterprise,
	}
	resp := &wallet.CreateUsersResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.WalletClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.CreateUsers(ctx, &req)
	if err != nil {
		log.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.CreateUsersRes{
		UserId: resp.UserId,
		Did:    resp.Did,
	}
	return result, nil
}

func (u user) UpdateUsers(ctx context.Context, params dto.UpdateUsers) (*dto.TxRes, error) {
	log := u.logger.WithContext(ctx).WithFields(
		map[string]interface{}{
			"function": "UpdateUsers",
			"params":   params,
		})

	req := wallet.UpdateUsersRequest{
		ProjectId: params.ProjectID,
		UserId:    params.UserId,
		PhoneNum:  params.PhoneNum,
	}
	resp := &wallet.UpdateUsersResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.WalletClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.UpdateUsers(ctx, &req)
	if err != nil {
		log.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.TxRes{}
	return result, nil
}

func (u user) ShowUsers(ctx context.Context, params dto.ShowUsers) (*dto.CreateUsersRes, error) {
	log := u.logger.WithContext(ctx).WithFields(
		map[string]interface{}{
			"function": "ShowUsers",
			"params":   params,
		})
	req := wallet.ShowUsersRequest{
		ProjectId: params.ProjectID,
		UserType:  wallet.USER_TYPE(params.Usertype),
		Code:      params.UserCode,
	}
	resp := &wallet.CreateUsersResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.WalletClientMap[mapKey]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err = grpcClient.ShowUsers(ctx, &req)
	if err != nil {
		log.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.CreateUsersRes{
		UserId: resp.UserId,
		Did:    resp.Did,
	}
	return result, nil
}
