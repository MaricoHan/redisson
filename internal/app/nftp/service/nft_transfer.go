package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (svc *NftTransfer) TransferNftClassByID(params dto.TransferNftClassByIDP) (string, error) {
	//不能自己转让给自己
	//400
	if params.Recipient == params.Owner {
		return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSelfTransfer)
	}

	//recipient不能为平台外账户或此应用外账户或非法账户
	_, err := models.TAccounts(
		models.TAccountWhere.AppID.EQ(params.AppID),
		models.TAccountWhere.Address.EQ(params.Recipient)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
		//400
		return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientFound)
	} else if err != nil {
		//500
		log.Error("transfer nft class", "query recipient error:", err.Error())
		return "", types.ErrInternal
	}

	//判断class
	class, err := models.TClasses(
		models.TClassWhere.ClassID.EQ(params.ClassID),
		models.TClassWhere.AppID.EQ(params.AppID),
		models.TClassWhere.Owner.EQ(params.Owner)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
		//404
		return "", types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("transfer nft class", "query class error:", err.Error())
		return "", types.ErrInternal
	}

	if class.Status != models.TClassesStatusActive {
		//400
		return "", types.ErrNftClassStatus
	}

	//msg
	msgs := nft.MsgTransferDenom{
		Id:        params.ClassID,
		Sender:    params.Owner,
		Recipient: params.Recipient,
	}

	//sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, _ := svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	baseTx.Gas = svc.base.calculateGas(data)
	data, hash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	if err != nil {
		log.Debug("transfer nft class", "BuildAndSign error:", err.Error())
		return "", types.ErrBuildAndSign
	}

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
		txId, err := svc.base.TxIntoDataBase(params.AppID,
			hash,
			data,
			models.TTXSOperationTypeTransferClass,
			models.TTXSStatusUndo, exec)
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
			return types.ErrInternal
		}

		return err
	})
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (svc *NftTransfer) TransferNftByIndex(params dto.TransferNftByIndexP) (string, error) {
	//不能自己转让给自己
	//400
	if params.Recipient == params.Owner {
		return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSelfTransfer)
	}

	//recipient不能为平台外账户或此应用外账户或非法账户
	_, err := models.TAccounts(
		models.TAccountWhere.AppID.EQ(params.AppID),
		models.TAccountWhere.Address.EQ(params.Recipient)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
		//400
		return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientFound)
	} else if err != nil {
		//500
		log.Error("transfer nft", "query recipient error:", err.Error())
		return "", types.ErrInternal
	}

	//msg
	res, err := models.TNFTS(models.TNFTWhere.Index.EQ(params.Index),
		models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
		models.TNFTWhere.AppID.EQ(params.AppID),
		models.TNFTWhere.Owner.EQ(params.Owner),
	).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
		//404
		return "", types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("transfer nft", "query nft error:", err.Error())
		return "", types.ErrInternal
	}

	//404
	if res.Status == models.TNFTSStatusBurned {
		return "", types.ErrNotFound
	}

	//400
	if res.Status != models.TNFTSStatusActive {
		return "", types.ErrNftStatus
	}

	msgs := nft.MsgTransferNFT{
		Id:        res.NFTID,
		DenomId:   string(params.ClassID),
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
	baseTx.Gas = svc.base.calculateGas(data)
	data, hash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	if err != nil {
		log.Debug("transfer nft by index", "BuildAndSign error:", err.Error())
		return "", types.ErrBuildAndSign
	}

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
		txId, err := svc.base.TxIntoDataBase(params.AppID,
			hash,
			data,
			models.TTXSOperationTypeTransferNFT,
			models.TTXSStatusUndo, exec)
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
		return "", err
	}

	return hash, nil
}

func (svc *NftTransfer) TransferNftByBatch(params dto.TransferNftByBatchP) (string, error) {
	indexMap := map[uint64]int{}
	var msgs sdktype.Msgs
	for i, modelResult := range params.Recipients {
		recipient := &dto.Recipient{
			Index:     modelResult.Index,
			Recipient: modelResult.Recipient,
		}

		//index校验 400
		if recipient.Index == 0 {
			return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(i)+"th "+types.ErrIndexInt)
		}

		//recipient不能为空 400
		if recipient.Recipient == "" {
			return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(i)+"th "+types.ErrRecipient)
		}

		//不能自己转让给自己
		//400
		if recipient.Recipient == params.Owner {
			return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(i)+"th "+types.ErrSelfTransfer)
		}

		//recipient不能为平台外账户或此应用外账户或非法账户
		_, err := models.TAccounts(
			models.TAccountWhere.AppID.EQ(params.AppID),
			models.TAccountWhere.Address.EQ(recipient.Recipient)).OneG(context.Background())
		if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
			(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
			//400
			return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(i)+"th "+types.ErrRecipientFound)
		} else if err != nil {
			//500
			log.Error("transfer nft by batch", "query recipient error:", err.Error())
			return "", types.ErrInternal
		}

		//判断index是否重复
		if _, ok := indexMap[recipient.Index]; ok {
			return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(i)+"th "+types.ErrRepeat)
		}

		res, err := models.TNFTS(models.TNFTWhere.Index.EQ(recipient.Index),
			models.TNFTWhere.ClassID.EQ(params.ClassID),
			models.TNFTWhere.AppID.EQ(params.AppID),
			models.TNFTWhere.Owner.EQ(params.Owner),
		).OneG(context.Background())
		if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
			(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
			//400
			return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(i)+"th "+types.ErrNftFound)
		} else if err != nil {
			//500
			log.Error("transfer nft by batch", "query nft error:", err.Error())
			return "", types.ErrInternal
		}

		//400
		if res.Status == models.TNFTSStatusBurned {
			return "", types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(i)+"th "+types.ErrNftFound)
		}

		//400
		if res.Status != models.TNFTSStatusActive {
			return "", types.NewAppError(types.RootCodeSpace, types.NftStatusAbnormal, "the "+string(i)+"th "+types.ErrNftStatusMsg)
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
		indexMap[recipient.Index] = 0
	}

	//sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, _ := svc.base.BuildAndSign(msgs, baseTx)
	baseTx.Gas = svc.base.calculateGas(data)
	data, hash, err := svc.base.BuildAndSign(msgs, baseTx)
	if err != nil {
		log.Debug("transfer nft by batch", "BuildAndSign error:", err.Error())
		return "", types.ErrBuildAndSign
	}

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
		txId, err := svc.base.TxIntoDataBase(params.AppID,
			hash,
			data,
			models.TTXSOperationTypeTransferNFTBatch,
			models.TTXSStatusUndo, exec)
		if err != nil {
			log.Debug("transfer nft by batch", "Tx Into DataBase error:", err.Error())
			return err
		}

		for j, modelResultR := range params.Recipients {
			recipient := &dto.Recipient{
				Index:     modelResultR.Index,
				Recipient: modelResultR.Recipient,
			}
			//index校验 400
			if recipient.Index == 0 {
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(j)+"th "+types.ErrIndexInt)
			}

			//recipient不能为空 400
			if recipient.Recipient == "" {
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(j)+"th "+types.ErrRecipient)
			}
			res, err := models.TNFTS(models.TNFTWhere.Index.EQ(recipient.Index),
				models.TNFTWhere.ClassID.EQ(params.ClassID),
				models.TNFTWhere.AppID.EQ(params.AppID),
				models.TNFTWhere.Owner.EQ(params.Owner),
			).One(context.Background(), exec)
			if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
				(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
				//400
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(j)+"th "+types.ErrNftFound)
			} else if err != nil {
				//500
				log.Error("transfer nft by batch", "query recipient error:", err.Error())
				return types.ErrInternal
			}

			if res.Status == models.TNFTSStatusBurned {
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+string(j)+"th "+types.ErrNftFound)
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
		return "", err
	}

	return hash, nil
}
