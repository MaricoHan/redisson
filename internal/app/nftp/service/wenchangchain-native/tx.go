package wenchangchain_native

import (
	"context"
	"database/sql"
	"strings"

	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/tendermint/tendermint/libs/json"
	"github.com/volatiletech/null/v8"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"

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
		log.Error("query tx by hash", "query tx error:", err.Error())
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

	if result.Status == 1 { //交易成功
		//根据 type 返回交易对象 id
		switch result.Type {
		case models.TTXSOperationTypeIssueClass:
			bytes := txinfo.Message.JSON
			var issueClass []nft.MsgIssueDenom
			err = json.Unmarshal(bytes, &issueClass)
			if err != nil {
				//500
				log.Error("query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = issueClass[0].Id
		case models.TTXSOperationTypeTransferClass:
			bytes := txinfo.Message.JSON
			var transferClass nft.MsgTransferDenom
			err = json.Unmarshal(bytes, &transferClass)
			if err != nil {
				//500
				log.Error("query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = transferClass.Id
		case models.TTXSOperationTypeMintNFT:
			bytes := txinfo.Message.JSON
			var mintNft []nft.MsgMintNFT
			err = json.Unmarshal(bytes, &mintNft)
			if err != nil {
				//500
				log.Error("query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = mintNft[0].DenomId

			nftTx, err := models.TMSGS(
				models.TDDCMSGWhere.TXHash.EQ(txinfo.Hash),
				models.TDDCMSGWhere.ProjectID.EQ(params.ProjectID),
			).OneG(context.Background())
			if err != nil {
				//500
				log.Error("query tx by hash", "query msg table error:", err.Error())
				return nil, types.ErrInternal
			}
			result.NftID = nftTx.NFTID.String
		case models.TTXSOperationTypeEditNFT:
			bytes := txinfo.Message.JSON
			var editNft nft.MsgEditNFT
			err = json.Unmarshal(bytes, &editNft)
			if err != nil {
				//500
				log.Error("query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = editNft.DenomId
			result.NftID = editNft.Id
		case models.TTXSOperationTypeBurnNFT:
			bytes := txinfo.Message.JSON
			var burnNft nft.MsgBurnNFT
			err = json.Unmarshal(bytes, &burnNft)
			if err != nil {
				//500
				log.Error("query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = burnNft.DenomId
			result.NftID = burnNft.Id
		case models.TTXSOperationTypeTransferNFT:
			bytes := txinfo.Message.JSON
			var transferNft nft.MsgTransferNFT
			err = json.Unmarshal(bytes, &transferNft)
			if err != nil {
				//500
				log.Error("query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = transferNft.DenomId
			result.NftID = transferNft.Id
		}
	}
	return result, nil
}
