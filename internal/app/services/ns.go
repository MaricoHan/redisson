package services

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	pb "gitlab.bianjie.ai/avata/chains/api/pb/v2/ns"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type INs interface {
	CreateDomain(ctx context.Context, params dto.CreateDomain) (*dto.TxRes, error)
	Domains(ctx context.Context, params dto.Domains) (*dto.DomainsRes, error)
}

type ns struct {
	logger *log.Entry
}

func NewNs(logger *log.Logger) *ns {
	return &ns{logger: logger.WithField("model", "ns")}
}

func (t *ns) CreateDomain(ctx context.Context, params dto.CreateDomain) (*dto.TxRes, error) {
	logger := t.logger.WithField("params", params).WithField("func", "CreateDomain")
	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.DomainCreateRequest{
		ProjectId:   params.ProjectID,
		OperationId: params.OperationId,
	}
	resp := &pb.DomainCreateResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NsClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Create(ctx, &req)
	if err != nil {
		logger.WithError(err).Error("request err")
		return nil, err
	}
	if err != nil {
		logger.Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.TxRes{}, nil
}

func (t *ns) Domains(ctx context.Context, params dto.Domains) (*dto.DomainsRes, error) {
	logger := t.logger.WithField("params", params).WithField("func", "Domains")
	// 非托管模式不支持
	if params.AccessMode == entity.UNMANAGED {
		return nil, errors2.ErrNotImplemented
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.DomainListRequest{
		ProjectId: params.ProjectID,
		Name:      params.Name,
		Tld:       params.Tld,
	}
	resp := &pb.DomainListResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.NsClientMap[mapKey]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.List(ctx, &req)
	if err != nil {
		logger.WithError(err).Error("request err")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	result := &dto.DomainsRes{
		Domains: []*dto.Domain{},
	}
	var domains []*dto.Domain
	for _, item := range resp.Data {
		domain := &dto.Domain{
			Name:            item.Name,
			Owner:           item.Owner,
			Status:          item.Status,
			Expire:          item.Expire,
			ExpireTimestamp: item.ExpireTimestamp,
		}
		domains = append(domains, domain)
	}
	if domains != nil {
		result.Domains = domains
	}

	return result, nil
}
