package services

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/types"
	pb "gitlab.bianjie.ai/avata/chains/api/pb/tx"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"time"
)

type ITx interface {
	TxResultByTxHash(params dto.TxResultByTxHash) (*dto.TxResultByTxHashRes, error)
}

type tx struct {
	logger *log.Logger
}

func NewTx(logger *log.Logger) *tx {
	return &tx{logger: logger}
}

func (t *tx) TxResultByTxHash(params dto.TxResultByTxHash) (*dto.TxResultByTxHashRes, error) {
	logFields := log.Fields{}
	logFields["model"] = "tx"
	logFields["func"] = "TxResultByTxHash"
	logFields["module"] = params.Module
	logFields["code"] = params.Code
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	req := pb.TxShowRequest{
		ProjectId: params.ProjectID,
		TaskId:    params.TaskId,
	}
	resp := &pb.TxShowResponse{}
	var err error
	mapKey := fmt.Sprintf("%s-%s", params.Code, params.Module)
	grpcClient, ok := initialize.TxClientMap[mapKey]
	if !ok {
		log.WithFields(logFields).Error(errors2.ErrService)
		return nil, errors2.New(errors2.InternalError, errors2.ErrService)
	}
	resp, err = grpcClient.Show(ctx, &req)
	if err != nil {
		log.WithFields(logFields).Error("request err:", err.Error())
		return nil, err
	}
	if resp == nil || resp.Detail == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.TxResultByTxHashRes{}
	status := pb.Status_value[resp.Detail.Status]
	result.Type = resp.Detail.OperationType
	result.TxHash = ""
	result.Status = status
	if status == int32(pb.Status_success) || status == int32(pb.Status_failed) {
		result.TxHash = resp.Detail.Hash
	}
	if resp.Detail.Tag != "" {
		var tagInterface interface{}
		err = json.Unmarshal([]byte(resp.Detail.Tag), &tagInterface)
		if err != nil {
			return nil, err
		}
		result.Tag = tagInterface.(map[string]interface{})
		if len(result.Tag) > 3 {
			return nil, constant.ErrInternal
		}
	}

	result.Message = resp.Detail.ErrMsg

	if result.Status == pb.Status_value["success"] { //交易成功
		//根据 type 返回交易对象 id
		result.BlockHeight = resp.Detail.BlockHeight
		result.Timestamp = resp.Detail.Timestamp
		typeJsonNft := types.JSON{}
		typeJsonMt := types.JSON{}
		if resp.Detail.Nft != "" {
			err = json.Unmarshal([]byte(resp.Detail.Nft), &typeJsonNft)
		}
		if resp.Detail.Mt != "" {
			err = json.Unmarshal([]byte(resp.Detail.Mt), &typeJsonMt)
		}
		if err != nil {
			return nil, err
		}
		result.Nft = typeJsonNft
		result.Mt = typeJsonMt
		//switch resp.Detail.OperationType {

		//case pb.OperationType_name[0]:
		//	result.ClassID = resp.Detail.ClassId
		//case pb.OperationType_name[5]:
		//	result.ClassID = resp.Detail.ClassId
		//default:
		//	result.ClassID = resp.Detail.ClassId
		//	result.NftID = resp.Detail.NftId
		//}
	}
	return result, nil
}
