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
