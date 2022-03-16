package wenchangchain_native

import (
	"context"
	"database/sql"
	"github.com/volatiletech/null/v8"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"strings"

	"github.com/friendsofgo/errors"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type Tx struct {
}

func NewTx() *service.TXBase {
	return &service.TXBase{
		Module:  service.NATIVE,
		Service: &Tx{},
	}
}

func (svc *Tx) Show(params dto.TxResultByTxHashP) (*dto.TxResultByTxHashRes, error) {
	//query
	txinfo, err := models.TTXS(
		models.TTXWhere.TaskID.EQ(null.StringFrom(params.TaskId)),
		models.TTXWhere.ProjectID.EQ(params.ProjectID),
	).OneG(context.Background())
	if err != nil {
		if (errors.Cause(err) == sql.ErrNoRows) || (strings.Contains(err.Error(), service.SqlNotFound)) {
			//404
			return nil, types.ErrNotFound
		}
		//500
		log.Error("ddc query tx by hash", "query tx error:", err.Error())
		return nil, types.ErrInternal
	}

	//result
	result := &dto.TxResultByTxHashRes{}
	result.Type = txinfo.OperationType
	result.TxHash = txinfo.Hash
	if txinfo.Status == models.TTXSStatusPending {
		result.Status = 0
	} else if txinfo.Status == models.TTXSStatusSuccess {
		result.Status = 1
	} else if txinfo.Status == models.TTXSStatusFailed {
		result.Status = 2
	} else {
		result.Status = 3 // tx.Status == "undo"
	}

	var tags map[string]interface{}
	err = txinfo.Tag.Unmarshal(&tags)
	if err != nil {
		//500
		log.Error("tx", "unmarshal error:", err.Error())
		return nil, types.ErrInternal
	}
	result.Message = txinfo.ErrMSG.String
	result.Tag = tags
	return result, nil
}
