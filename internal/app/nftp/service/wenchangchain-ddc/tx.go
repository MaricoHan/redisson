package wenchangchain_ddc

import (
	"context"
	"encoding/json"

	"strings"

	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"

	"github.com/irisnet/irismod-sdk-go/nft"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type Tx struct {
}

func NewTx() *service.TXBase {
	return &service.TXBase{
		Module:  service.DDC,
		Service: &Tx{},
	}
}

func (t *Tx) Show(params dto.TxResultByTxHashP) (*dto.TxResultByTxHashRes, error) {
	//query
	txinfo, err := models.TDDCTXS(
		models.TDDCTXWhere.TaskID.EQ(null.StringFrom(params.TaskId)),
		models.TDDCTXWhere.ProjectID.EQ(params.ProjectID),
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
	if txinfo.Status == models.TDDCTXSStatusPending {
		result.Status = 0
	} else if txinfo.Status == models.TDDCTXSStatusSuccess {
		result.Status = 1
	} else if txinfo.Status == models.TDDCTXSStatusFailed {
		result.Status = 2
	} else {
		result.Status = 3 // tx.Status == "undo"
	}

	var tags map[string]interface{}
	err = txinfo.Tag.Unmarshal(&tags)
	if err != nil {
		//500
		log.Error("ddc tx", "unmarshal error:", err.Error())
		return nil, types.ErrInternal
	}
	result.Message = txinfo.ErrMSG.String
	result.Tag = tags

	if result.Status == 1 { //交易成功
		//根据 type 返回交易对象 id
		switch result.Type {
		case models.TDDCTXSOperationTypeIssueClass:
			bytes := txinfo.Message.JSON
			var issueClass []nft.MsgIssueDenom
			err = json.Unmarshal(bytes, &issueClass)
			if err != nil {
				//500
				log.Error("ddc query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = issueClass[0].Id
		case models.TDDCTXSOperationTypeTransferClass:
			bytes := txinfo.Message.JSON
			var transferClass []nft.MsgTransferDenom
			err = json.Unmarshal(bytes, &transferClass)
			if err != nil {
				//500
				log.Error("ddc query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = transferClass[0].Id
		case models.TDDCTXSOperationTypeMintNFT:
			bytes := txinfo.Message.JSON
			var mintNft []nft.MsgMintNFT
			err = json.Unmarshal(bytes, &mintNft)
			if err != nil {
				//500
				log.Error("ddc query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = mintNft[0].DenomId

			ddcTx, err := models.TDDCMSGS(
				models.TDDCMSGWhere.TXHash.EQ(txinfo.Hash),
				models.TDDCMSGWhere.ProjectID.EQ(params.ProjectID),
			).OneG(context.Background())
			if err != nil {
				//500
				log.Error("ddc query tx by hash", "query msg table error:", err.Error())
				return nil, types.ErrInternal
			}
			result.NftID = ddcTx.NFTID.String
		case models.TDDCTXSOperationTypeEditNFT:
			bytes := txinfo.Message.JSON
			var editNft []nft.MsgEditNFT
			err = json.Unmarshal(bytes, &editNft)
			if err != nil {
				//500
				log.Error("ddc query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = editNft[0].DenomId
			result.NftID = editNft[0].Id
		case models.TDDCTXSOperationTypeBurnNFT:
			bytes := txinfo.Message.JSON
			var burnNft []nft.MsgBurnNFT
			err = json.Unmarshal(bytes, &burnNft)
			if err != nil {
				//500
				log.Error("ddc query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = burnNft[0].DenomId
			result.NftID = burnNft[0].Id
		case models.TDDCTXSOperationTypeTransferNFT:
			bytes := txinfo.Message.JSON
			var transferNft []nft.MsgTransferNFT
			err = json.Unmarshal(bytes, &transferNft)
			if err != nil {
				//500
				log.Error("ddc query tx by hash", "unmarshal tx message error:", err.Error())
				return nil, types.ErrInternal
			}
			result.ClassID = transferNft[0].DenomId
			result.NftID = transferNft[0].Id
		}
	}
	return result, nil
}
