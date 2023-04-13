package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	records "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/record"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
)

type IRecord interface {
	CreateRecord(ctx context.Context, params *records.RecordCreateRequest) (*records.RecordCreateResponse, error)
}
type record struct {
	logger *log.Logger
}

func NewRecord(logger *log.Logger) *record {
	return &record{logger: logger}
}

// CreateRecord 创建存证
func (r *record) CreateRecord(ctx context.Context, params *records.RecordCreateRequest) (*records.RecordCreateResponse, error) {
	authData := r.authData(ctx)
	params.ProjectId = authData.ProjectId
	mapKey := fmt.Sprintf("%s-%s", authData.Code, authData.Module)
	logger := r.logger.WithContext(ctx).WithFields(map[string]interface{}{
		"server-name": mapKey,
		"params":      params,
		"func":        "CreateRecord",
	})
	grpcClient, ok := initialize.RecordClientMap[mapKey]
	if !ok {
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.Create(ctx, params)
	if err != nil {
		logger.WithError(err).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return resp, nil
}

func (r *record) authData(ctx context.Context) vo.AuthData {
	authDataString := ctx.Value("X-Auth-Data")
	authDataSlice, ok := authDataString.([]string)
	if !ok {
		return vo.AuthData{}
	}
	var authData vo.AuthData
	err := json.Unmarshal([]byte(authDataSlice[0]), &authData)
	if err != nil {
		log.Error("auth data Error: ", err)
		return vo.AuthData{}
	}
	return authData
}
