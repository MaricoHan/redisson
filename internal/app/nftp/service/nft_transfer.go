package service

import (
	"context"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
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
		return "", types.ErrBuildAndSign
	}

	//txs status = undo
	txId, err := svc.base.TxIntoDataBase(params.AppID,
		hash,
		data,
		models.TTXSOperationTypeTransferClass,
		models.TTXSStatusUndo)
	if err != nil {
		return "", err
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//class status && class.lockby == txid
		class, err := models.TClasses(
			models.TClassWhere.ClassID.EQ(string(params.ClassID)),
			models.TClassWhere.AppID.EQ(params.AppID),
			models.TClassWhere.Owner.EQ(params.Owner)).One(context.Background(), exec)
		if err != nil {
			return types.ErrNftClassTransfer
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
	//msg
	res, err := models.TNFTS(models.TNFTWhere.Index.EQ(params.Index),
		models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
		models.TNFTWhere.AppID.EQ(params.AppID),
		models.TNFTWhere.Owner.EQ(params.Owner),
	).OneG(context.Background())
	if err != nil {
		return "", types.ErrNftTransfer
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
		return "", types.ErrBuildAndSign
	}

	//写入txs status = undo
	txId, err := svc.base.TxIntoDataBase(params.AppID,
		hash,
		data,
		models.TTXSOperationTypeTransferNFT,
		models.TTXSStatusUndo)
	if err != nil {
		return "", err
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
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
		res, err := models.TNFTS(models.TNFTWhere.Index.EQ(recipient.Index),
			models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
			models.TNFTWhere.AppID.EQ(params.AppID),
			models.TNFTWhere.Owner.EQ(params.Owner),
		).OneG(context.Background())
		if err != nil {
			return "", types.ErrNftBatchTransfer
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
		return "", types.ErrBuildAndSign
	}

	//写入txs status = undo
	txId, err := svc.base.TxIntoDataBase(params.AppID,
		hash,
		data,
		models.TTXSOperationTypeTransferNFTBatch,
		models.TTXSStatusUndo)
	if err != nil {
		return "", err
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
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
		return "", err
	}

	return hash, nil
}
