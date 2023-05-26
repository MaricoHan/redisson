package native

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/dict"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"

	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
)

type IDict interface {
	ListTxTypes(ctx context.Context, params *dto.ListTxTypes) (*vo.TxTypesRes, error)
}

type dict struct {
	logger *log.Entry
}

var _ IDict = dict{}

func NewDict(logger *log.Logger) *dict {
	return &dict{logger.WithField("model", "dict")}
}

func (d dict) ListTxTypes(ctx context.Context, params *dto.ListTxTypes) (*vo.TxTypesRes, error) {
	log := d.logger.WithContext(ctx).WithField("func", "ListTxTypes")

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	grpcClient, ok := initialize.NativeDictClientMap[fmt.Sprintf("%s-%s", params.Code, params.Module)]
	if !ok {
		log.Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}

	resp, err := grpcClient.ListTxTypes(ctx, &pb.ListTxTypesReq{})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	res := &vo.TxTypesRes{}
	for _, data := range resp.Data {
		res.Data = append(res.Data, &vo.TxType{
			Module:      data.Module,
			Operation:   data.Operation,
			Code:        data.Code,
			Name:        data.Name,
			Description: data.Description,
		})
	}
	return res, nil
}
