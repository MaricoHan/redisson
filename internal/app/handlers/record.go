package handlers

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/record"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"
	errors2 "gitlab.bianjie.ai/avata/utils/errors/v2"

	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
)

type IRecord interface {
	CreateRecord(ctx context.Context, _ interface{}) (interface{}, error)
}

type Record struct {
	base.Base
	base.PageBasic
	svc services.IRecord
}

func NewRecord(svc services.IRecord) *Record {
	return &Record{svc: svc}
}

func (r *Record) CreateRecord(ctx context.Context, request interface{}) (interface{}, error) {
	recordParams := request.(*record.RecordCreateRequest)
	operationId := strings.TrimSpace(recordParams.OperationId)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	return r.svc.CreateRecord(ctx, recordParams)
}
