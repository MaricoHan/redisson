package service

import (
	"context"

	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type NftTransfer struct {
	base *Base
}

func NewNftTransfer(base *Base) *NftTransfer {
	return &NftTransfer{base: base}
}

func (svc *NftTransfer) TransferNftClassByID(params dto.TransferNftClassByIDP) (string, error) {
	db, err := orm.GetDB().Begin()
	if err != nil {
		return "", types.ErrMysqlConn
	}

	//query if class can be operated
	class, err := models.TClasses(
		models.TClassWhere.ClassID.EQ(string(params.ClassID)),
		models.TClassWhere.AppID.EQ(params.AppID),
		models.TClassWhere.Owner.EQ(params.Owner)).OneG(context.Background())
	if err != nil {
		db.Rollback()
		return "", types.ErrNftClassTransfer
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
		return "", types.ErrBuildAndSign
	}

	//txs status = undo
	txs := models.TTX{
		AppID:         params.AppID,
		Hash:          hash,
		OriginData:    null.BytesFrom(data),
		OperationType: "transfer_class",
		Status:        "undo",
	}
	err = txs.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftClassTransfer
	}

	//class status = pending && lockby = txs.id
	class.Status = "pending"
	class.LockedBy = txs.ID
	_, err = class.UpdateG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftClassTransfer
	}

	err = db.Commit()
	if err != nil {
		return "", types.ErrInternal
	}

	return hash, nil
}

func (svc *NftTransfer) TransferNftByIndex(params dto.TransferNftByIndexP) (string, error) {
	db, err := orm.GetDB().Begin()
	if err != nil {
		return "", types.ErrMysqlConn
	}

	//msg
	res, err := models.TNFTS(models.TNFTWhere.Index.EQ(params.Index),
		models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
		models.TNFTWhere.AppID.EQ(params.AppID),
		models.TNFTWhere.Owner.EQ(params.Owner),
	).OneG(context.Background())
	if err != nil {
		db.Rollback()
		return "", types.ErrNftTransfer
	}

	msgs := nft.MsgTransferNFT{
		Id:        res.NFTID,
		DenomId:   string(params.ClassID),
		Sender:    params.Owner,
		Recipient: params.Recipient,
	}

	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	if err != nil {
		return "", types.ErrBuildAndSign
	}
	//写入txs status = undo
	txs := models.TTX{
		AppID:         params.AppID,
		Hash:          hash,
		OriginData:    null.BytesFrom(data),
		OperationType: "transfer_nft",
		Status:        "undo",
	}
	err = txs.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftTransfer
	}

	res.Status = "pending"
	res.LockedBy = null.Uint64From(txs.ID)
	_, err = res.UpdateG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftTransfer
	}

	err = db.Commit()
	if err != nil {
		return "", types.ErrInternal
	}

	return hash, nil
}

func (svc *NftTransfer) TransferNftByBatch(params dto.TransferNftByBatchP) (string, error) {
	db, err := orm.GetDB().Begin()
	if err != nil {
		return "", types.ErrMysqlConn
	}

	var msgs sdktype.Msgs
	for _, modelResult := range params.Recipients {
		recipient := &dto.Recipient{
			Index:     modelResult.Index,
			Recipient: modelResult.Recipient,
		}
		//msg
		res, err := models.TNFTS(models.TNFTWhere.Index.EQ(recipient.Index),
			models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
			models.TNFTWhere.AppID.EQ(params.AppID),
			models.TNFTWhere.Owner.EQ(params.Owner),
		).OneG(context.Background())
		if err != nil {
			db.Rollback()
			return "", types.ErrNftBatchTransfer
		}
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
	txs := models.TTX{
		AppID:         params.AppID,
		Hash:          hash,
		OriginData:    null.BytesFrom(data),
		OperationType: "transfer_nft_batch",
		Status:        "undo",
	}
	err = txs.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftBatchTransfer
	}

	for _, modelResultR := range params.Recipients {
		recipient := &dto.Recipient{
			Index:     modelResultR.Index,
			Recipient: modelResultR.Recipient,
		}
		res, err := models.TNFTS(models.TNFTWhere.Index.EQ(recipient.Index),
			models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
			models.TNFTWhere.AppID.EQ(params.AppID),
			models.TNFTWhere.Owner.EQ(params.Owner),
		).OneG(context.Background())
		if err != nil {
			db.Rollback()
			return "", types.ErrNftBatchTransfer
		}

		res.Status = "pending"
		res.LockedBy = null.Uint64From(txs.ID)
		_, err = res.UpdateG(context.Background(), boil.Infer())
		if err != nil {
			return "", types.ErrNftBatchTransfer
		}
	}

	err = db.Commit()
	if err != nil {
		return "", types.ErrInternal
	}
	return hash, nil
}
