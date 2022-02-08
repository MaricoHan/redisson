package service

import (
	"context"
	"database/sql"
	"github.com/friendsofgo/errors"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type Tx struct {
}

func NewTx() *Tx {
	return &Tx{}
}

func (svc *Tx) TxResultByTxHash(params dto.TxResultByTxHashP) (*dto.TxResultByTxHashRes, error) {
	//query
	txinfo, err := models.TTXS(
		models.TTXWhere.Hash.EQ(params.Hash),
		models.TTXWhere.AppID.EQ(params.AppID),
	).OneG(context.Background())
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		//404
		return nil, types.ErrTxNotFound
	} else if err != nil {
		//500
		log.Error("query tx by hash", "query tx error:", err.Error())
		return nil, types.ErrQuery
	}

	//result
	result := &dto.TxResultByTxHashRes{}
	result.Type = txinfo.OperationType

	if txinfo.Status == models.TTXSStatusPending {
		result.Status = 0
	} else if txinfo.Status == models.TTXSStatusSuccess {
		result.Status = 1
	} else {
		result.Status = 2 // tx.Status == "failed"
	}
	result.Message = txinfo.ErrMSG.String

	return result, nil
}
