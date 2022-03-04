package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
)

type NftTransfer struct {
	base *Base
}

func NewNftTransfer(base *Base) *NftTransfer {
	return &NftTransfer{base: base}
}

func (svc *NftTransfer) TransferNftClassByID(params dto.TransferNftClassByIDP) (*dto.TxRes, error) {
	//不能自己转让给自己
	//400
	if params.Recipient == params.Owner {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSelfTransfer)
	}

	// ValidateSigner
	if err := svc.base.ValidateSigner(params.Owner, params.ProjectID); err != nil {
		return nil, err
	}

	// ValidateRecipient
	if err := svc.base.ValidateRecipient(params.Recipient, params.ProjectID); err != nil {
		return nil, err
	}
	//判断class
	class, err := models.TClasses(
		models.TClassWhere.ClassID.EQ(params.ClassID),
		models.TClassWhere.ProjectID.EQ(params.ProjectID),
		models.TClassWhere.Owner.EQ(params.Owner)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
		//404
		return nil, types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("transfer nft class", "query class error:", err.Error())
		return nil, types.ErrInternal
	}

	if class.Status != models.TClassesStatusActive {
		//400
		return nil, types.ErrNftClassStatus
	}

	//msg
	msgs := nft.MsgTransferDenom{
		Id:        params.ClassID,
		Sender:    params.Owner,
		Recipient: params.Recipient,
	}

	//sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	baseTx.Gas = svc.base.transferDenomGas(class)
	err = svc.base.GasThan(params.Owner, params.ChainID, baseTx.Gas, params.PlatFormID)
	if err != nil {
		return nil, types.ErrInternal
	}
	data, hash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	if err != nil {
		log.Debug("transfer nft class", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//validate tx
		txone, err := svc.base.ValidateTx(hash)
		if err != nil {
			return err
		}
		if txone != nil && txone.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txone.ID)
			data, hash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
			if err != nil {
				log.Debug("transfer nft class", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		//txs status = undo
		messageByte, _ := json.Marshal(msgs)
		code := fmt.Sprintf("%s%s%s", params.Owner, models.TTXSOperationTypeTransferClass, time.Now().String())
		taskId = svc.base.EncodeData(code)
		// Tx into database
		txId, err := svc.base.UndoTxIntoDataBase(params.Owner, models.TTXSOperationTypeTransferClass, taskId, hash,
			params.ProjectID, data, messageByte, params.Tag, int64(baseTx.Gas), exec)

		if err != nil {
			log.Debug("transfer nft class", "Tx Into DataBase error:", err.Error())
			return err
		}

		//class status = pending && lockby = txs.id
		class.Status = models.TClassesStatusPending
		class.LockedBy = null.Uint64From(txId)

		ok, err := class.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			//500
			log.Error("transfer nft class", "update class error:", err.Error())
			return types.ErrInternal
		}
		if ok != 1 {
			log.Error("transfer nft class", "update class error:", err.Error())
			return types.ErrInternal
		}

		return err
	})
	if err != nil {
		return nil, err
	}
	return &dto.TxRes{TaskId: taskId}, nil
}

func (svc *NftTransfer) TransferNftByNftId(params dto.TransferNftByNftIdP) (*dto.TxRes, error) {
	//不能自己转让给自己
	//400
	if params.Recipient == params.Owner {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSelfTransfer)
	}
	// ValidateSigner
	if err := svc.base.ValidateSigner(params.Owner, params.ProjectID); err != nil {
		return nil, err
	}

	// ValidateRecipient
	if err := svc.base.ValidateRecipient(params.Recipient, params.ProjectID); err != nil {
		return nil, err
	}

	//msg
	res, err := models.TNFTS(
		models.TNFTWhere.NFTID.EQ(params.NftId),
		models.TNFTWhere.ClassID.EQ(params.ClassID),
		models.TNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TNFTWhere.Owner.EQ(params.Owner),
	).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
		//404
		return nil, types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("transfer nft", "query nft error:", err.Error())
		return nil, types.ErrInternal
	}

	//404
	if res.Status == models.TNFTSStatusBurned {
		return nil, types.ErrNotFound
	}

	//400
	if res.Status != models.TNFTSStatusActive {
		return nil, types.ErrNftStatus
	}

	msgs := nft.MsgTransferNFT{
		Id:        res.NFTID,
		DenomId:   params.ClassID,
		Name:      res.Name.String,
		URI:       res.URI.String,
		Data:      res.Metadata.String,
		Sender:    params.Owner,
		Recipient: params.Recipient,
		UriHash:   res.URIHash.String,
	}

	//build and sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, _ := svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	baseTx.Gas = svc.base.transferOneNftGas(data)
	err = svc.base.GasThan(params.Owner, params.ChainID, baseTx.Gas, params.PlatFormID)
	if err != nil {
		return nil, types.ErrInternal
	}
	data, hash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	if err != nil {
		log.Debug("transfer nft by index", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//validate tx
		txone, err := svc.base.ValidateTx(hash)
		if err != nil {
			return err
		}
		if txone != nil && txone.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txone.ID)
			data, hash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
			if err != nil {
				log.Debug("transfer nft by index", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		//写入txs status = undo
		messageByte, _ := json.Marshal(msgs)
		code := fmt.Sprintf("%s%s%s", params.Owner, models.TTXSOperationTypeTransferNFT, time.Now().String())
		taskId = svc.base.EncodeData(code)

		// Tx into database
		txId, err := svc.base.UndoTxIntoDataBase(params.Owner, models.TTXSOperationTypeTransferNFT, taskId, hash,
			params.ProjectID, data, messageByte, params.Tag, int64(baseTx.Gas), exec)

		if err != nil {
			log.Debug("transfer nft by index", "Tx Into DataBase error:", err.Error())
			return types.ErrInternal
		}

		res.Status = models.TNFTSStatusPending
		res.LockedBy = null.Uint64From(txId)
		ok, err := res.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			return types.ErrInternal
		}
		if ok != 1 {
			return types.ErrInternal
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return &dto.TxRes{TaskId: taskId}, nil
}

func (svc *NftTransfer) TransferNftByBatch(params dto.TransferNftByBatchP) (*dto.TxRes, error) {
	nftIdMap := map[string]int{}
	var msgs sdktype.Msgs
	var amount uint64
	for i, modelResult := range params.Recipients {
		recipient := &dto.Recipient{
			NftId:     modelResult.NftId,
			Recipient: modelResult.Recipient,
		}
		recipient.Recipient = strings.TrimSpace(recipient.Recipient)

		//index校验 400
		if recipient.NftId == "" {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftIdLen+" or "+types.ErrNftIdString)
		}

		//recipient不能为空 400
		if recipient.Recipient == "" {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrRecipient)
		}

		if len([]rune(recipient.Recipient)) > 128 {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrRecipientLen)
		}

		//不能自己转让给自己
		//400
		if recipient.Recipient == params.Owner {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrSelfTransfer)
		}

		//判断nft_id是否重复
		if _, ok := nftIdMap[recipient.NftId]; ok {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrRepeat)
		}

		res, err := models.TNFTS(models.TNFTWhere.NFTID.EQ(recipient.NftId),
			models.TNFTWhere.ClassID.EQ(params.ClassID),
			models.TNFTWhere.ProjectID.EQ(params.ProjectID),
			models.TNFTWhere.Owner.EQ(params.Owner),
		).OneG(context.Background())
		if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
			(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
			//400
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftFound)
		} else if err != nil {
			//500
			log.Error("transfer nft by batch", "query nft error:", err.Error())
			return nil, types.ErrInternal
		}

		//400
		if res.Status == models.TNFTSStatusBurned {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftFound)
		}

		//400
		if res.Status != models.TNFTSStatusActive {
			return nil, types.NewAppError(types.RootCodeSpace, types.NftStatusAbnormal, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftStatusMsg)
		}

		//msg
		msg := nft.MsgTransferNFT{
			Id:        res.NFTID,
			DenomId:   params.ClassID,
			Name:      res.Name.String,
			URI:       res.URI.String,
			Data:      res.Metadata.String,
			Sender:    params.Owner,
			Recipient: recipient.Recipient,
			UriHash:   res.URIHash.String,
		}
		msgs = append(msgs, &msg)
		nftIdMap[recipient.NftId] = 0
		amount += 1
	}

	//sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, _ := svc.base.BuildAndSign(msgs, baseTx)
	baseTx.Gas = svc.base.transferNftsGas(data, amount)
	data, hash, err := svc.base.BuildAndSign(msgs, baseTx)
	if err != nil {
		log.Debug("transfer nft by batch", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//validate tx
		txone, err := svc.base.ValidateTx(hash)
		if err != nil {
			return err
		}
		if txone != nil && txone.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txone.ID)
			data, hash, err = svc.base.BuildAndSign(msgs, baseTx)
			if err != nil {
				log.Debug("transfer nft by batch", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		//写入txs status = undo
		messageByte, _ := json.Marshal(msgs)
		code := fmt.Sprintf("%s%s%s", params.Owner, models.TTXSOperationTypeTransferNFTBatch, time.Now().String())
		taskId = svc.base.EncodeData(code)

		// Tx into database
		txId, err := svc.base.UndoTxIntoDataBase(params.Owner, models.TTXSOperationTypeTransferNFTBatch, taskId, hash,
			params.ProjectID, data, messageByte, nil, int64(baseTx.Gas), exec)

		if err != nil {
			log.Debug("transfer nft by batch", "Tx Into DataBase error:", err.Error())
			return err
		}

		for j, modelResultR := range params.Recipients {
			recipient := &dto.Recipient{
				NftId:     modelResultR.NftId,
				Recipient: modelResultR.Recipient,
			}
			recipient.Recipient = strings.TrimSpace(recipient.Recipient)
			//index校验 400
			if recipient.NftId == "" {
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", j+1)+"th "+types.ErrNftIdString)
			}

			//recipient不能为空 400
			if recipient.Recipient == "" {
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", j+1)+types.ErrRecipient)
			}
			res, err := models.TNFTS(models.TNFTWhere.NFTID.EQ(recipient.NftId),
				models.TNFTWhere.ClassID.EQ(params.ClassID),
				models.TNFTWhere.ProjectID.EQ(params.ProjectID),
				models.TNFTWhere.Owner.EQ(params.Owner),
			).One(context.Background(), exec)
			if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
				(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
				//400
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", j+1)+types.ErrNftFound)
			} else if err != nil {
				//500
				log.Error("transfer nft by batch", "query recipient error:", err.Error())
				return types.ErrInternal
			}

			if res.Status == models.TNFTSStatusBurned {
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", j+1)+types.ErrNftFound)
			}

			if res.Status != models.TNFTSStatusActive {
				return types.ErrNftStatus
			}

			res.Status = models.TNFTSStatusPending
			res.LockedBy = null.Uint64From(txId)
			ok, err := res.Update(context.Background(), exec, boil.Infer())
			if err != nil {
				return types.ErrInternal
			}
			if ok != 1 {
				return types.ErrInternal
			}
		}
		return err
	})
	if err != nil {
		//自定义err
		return nil, err
	}
	return &dto.TxRes{TaskId: taskId}, nil
}
