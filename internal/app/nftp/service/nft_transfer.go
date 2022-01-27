package service

import (
	"context"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
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
	acc, err := models.TAccounts(
		models.TAccountWhere.AppID.EQ(params.AppID),
		models.TAccountWhere.Address.EQ(params.Recipient)).OneG(context.Background())
	if err != nil {
		return "", types.ErrParams
	}
	if acc == nil {
		return "", types.ErrParams
	}

	class, err := models.TClasses(
		models.TClassWhere.ClassID.EQ(string(params.ClassID)),
		models.TClassWhere.AppID.EQ(params.AppID),
		models.TClassWhere.Owner.EQ(params.Owner)).OneG(context.Background())
	if err != nil {
		return "", types.ErrNftClassTransfer
	}

	if class.Status != models.TClassesStatusActive {
		return "", types.ErrNftClassStatus
	}

	//msg
	msgs := nft.MsgTransferDenom{
		Id:        string(params.ClassID),
		Sender:    params.Owner,
		Recipient: params.Recipient,
	}

	//sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)

	if err != nil {
		log.Debug("transfer nft class", "BuildAndSign error:", err.Error())
		return "", types.ErrBuildAndSign
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
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
			return types.ErrNftClassTransfer
		}
		if ok != 1 {
			return types.ErrNftClassTransfer
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (svc *NftTransfer) TransferNftByIndex(params dto.TransferNftByIndexP) (string, error) {
	acc, err := models.TAccounts(
		models.TAccountWhere.AppID.EQ(params.AppID),
		models.TAccountWhere.Address.EQ(params.Recipient)).OneG(context.Background())
	if err != nil {
		return "", types.ErrParams
	}
	if acc == nil {
		return "", types.ErrParams
	}

	//msg
	res, err := models.TNFTS(models.TNFTWhere.Index.EQ(params.Index),
		models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
		models.TNFTWhere.AppID.EQ(params.AppID),
		models.TNFTWhere.Owner.EQ(params.Owner),
	).OneG(context.Background())
	if err != nil {
		return "", types.ErrNftTransfer
	}

	if res.Status != models.TNFTSStatusActive {
		return "", types.ErrNftStatus
	}

	msgs := nft.MsgTransferNFT{
		Id:        res.NFTID,
		DenomId:   string(params.ClassID),
		Sender:    params.Owner,
		Recipient: params.Recipient,
	}

	//build and sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	if err != nil {
		log.Debug("transfer nft by index", "BuildAndSign error:", err.Error())
		return "", types.ErrBuildAndSign
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//写入txs status = undo
		txId, err := svc.base.TxIntoDataBase(params.AppID,
			hash,
			data,
			models.TTXSOperationTypeTransferNFT,
			models.TTXSStatusUndo, exec)
		if err != nil {
			log.Debug("transfer nft by index", "Tx Into DataBase error:", err.Error())
			return types.ErrTxMsgInsert
		}

		res.Status = models.TNFTSStatusPending
		res.LockedBy = null.Uint64From(txId)
		ok, err := res.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			return types.ErrNftTransfer
		}
		if ok != 1 {
			return types.ErrNftTransfer
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (svc *NftTransfer) TransferNftByBatch(params dto.TransferNftByBatchP) (string, error) {
	var msgs sdktype.Msgs
	for _, modelResult := range params.Recipients {
		recipient := &dto.Recipient{
			Index:     modelResult.Index,
			Recipient: modelResult.Recipient,
		}
		if recipient.Index == 0 {
			return "", types.ErrParams
		}
		if recipient.Recipient == "" {
			return "", types.ErrParams
		}
		acc, err := models.TAccounts(
			models.TAccountWhere.AppID.EQ(params.AppID),
			models.TAccountWhere.Address.EQ(recipient.Recipient)).OneG(context.Background())
		if err != nil {
			return "", types.ErrParams
		}
		if acc == nil {
			return "", types.ErrParams
		}

		res, err := models.TNFTS(models.TNFTWhere.Index.EQ(recipient.Index),
			models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
			models.TNFTWhere.AppID.EQ(params.AppID),
			models.TNFTWhere.Owner.EQ(params.Owner),
		).OneG(context.Background())
		if err != nil {
			return "", types.ErrNftBatchTransfer
		}

		if res.Status != models.TNFTSStatusActive {
			return "", types.ErrNftStatus
		}

		//msg
		msg := nft.MsgTransferNFT{
			Id:        res.NFTID,
			DenomId:   string(params.ClassID),
			Sender:    params.Owner,
			Recipient: recipient.Recipient,
		}
		msgs = append(sdktype.Msgs{&msg})
	}

	//sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, err := svc.base.BuildAndSign(msgs, baseTx)
	if err != nil {
		log.Debug("transfer nft by batch", "BuildAndSign error:", err.Error())
		return "", types.ErrBuildAndSign
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
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

		for _, modelResultR := range params.Recipients {
			recipient := &dto.Recipient{
				Index:     modelResultR.Index,
				Recipient: modelResultR.Recipient,
			}
			if recipient.Index == 0 {
				return types.ErrParams
			}
			if recipient.Recipient == "" {
				return types.ErrParams
			}
			res, err := models.TNFTS(models.TNFTWhere.Index.EQ(recipient.Index),
				models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
				models.TNFTWhere.AppID.EQ(params.AppID),
				models.TNFTWhere.Owner.EQ(params.Owner),
			).One(context.Background(), exec)
			if err != nil {
				return types.ErrNftBatchTransfer
			}

			if res.Status != models.TNFTSStatusActive {
				return types.ErrNftStatus
			}

			res.Status = models.TNFTSStatusPending
			res.LockedBy = null.Uint64From(txId)
			ok, err := res.Update(context.Background(), exec, boil.Infer())
			if err != nil {
				return types.ErrNftBatchTransfer
			}
			if ok != 1 {
				return types.ErrNftClassesSet
			}
		}
		return nil
	})
	if err != nil {
		//自定义err
		return "", err
	}

	return hash, nil
}
