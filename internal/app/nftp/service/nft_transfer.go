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
		return "", types.ErrTxResult
	}

	//txs status = undo
	txs := models.TTX{
		AppID:      params.AppID,
		Hash:       hash,
		Status:     "undo",
		OriginData: null.BytesFrom(data),
	}
	err = txs.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftClassTransfer
	}

	//class status = pendding && lockby = txs.id
	class, err := models.TClasses(
		models.TClassWhere.ClassID.EQ(string(params.ClassID)),
		models.TClassWhere.AppID.EQ(params.AppID),
		models.TClassWhere.Owner.EQ(params.Owner)).OneG(context.Background())
	if err != nil {
		return "", types.ErrTxResult
	}
	class.Status = "pendding"

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
		return "", types.ErrMysqlConn
	}
	msgs := nft.MsgTransferNFT{
		Id:        res.NFTID,
		DenomId:   string(params.ClassID),
		Sender:    params.Owner,
		Recipient: params.Recipient,
	}

	//sign
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	data, hash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgs}, baseTx)
	if err != nil {
		return "", types.ErrTxResult
	}

	//写入txs status = undo
	txs := models.TTX{
		AppID:      params.AppID,
		Hash:       hash,
		Status:     "undo",
		OriginData: null.BytesFrom(data),
	}
	err = txs.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftClassTransfer
	}

	//nft status = pendding && lockby = txs.id
	nft, err := models.TNFTS(
		models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
		models.TNFTWhere.ClassID.EQ(res.NFTID),
		models.TNFTWhere.AppID.EQ(params.AppID),
		models.TNFTWhere.Owner.EQ(params.Owner),
	).OneG(context.Background())
	if err != nil {
		return "", types.ErrNftClassTransfer
	}

	nft.Status = "pendding"
	_, err = nft.UpdateG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftClassTransfer
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
			return "", types.ErrNftClassTransfer
		}
		msg := nft.MsgTransferNFT{
			Id:        res.NFTID,
			DenomId:   string(params.ClassID),
			Sender:    params.Owner,
			Recipient: recipient.Recipient,
		}
		msgs = append(sdktype.Msgs{&msg})

		//nft status = pendding && lockby = txs.id
		nft, err := models.TNFTS(
			models.TNFTWhere.ClassID.EQ(string(params.ClassID)),
			models.TNFTWhere.ClassID.EQ(res.NFTID),
			models.TNFTWhere.AppID.EQ(params.AppID),
			models.TNFTWhere.Owner.EQ(params.Owner),
		).OneG(context.Background())
		if err != nil {
			return "", types.ErrTxResult
		}
		nft.Status = "pendding"
		_, err = nft.UpdateG(context.Background(), boil.Infer())
		if err != nil {
			return "", types.ErrNftClassTransfer
		}
	}

	//sign
	baseTx := svc.base.CreateBaseTx("", "")
	data, hash, err := svc.base.BuildAndSign(msgs, baseTx)
	if err != nil {
		return "", types.ErrTxResult
	}

	//写入txs status = undo
	txs := models.TTX{
		AppID:      params.AppID,
		Hash:       hash,
		Status:     "undo",
		OriginData: null.BytesFrom(data),
	}
	err = txs.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return "", types.ErrNftClassTransfer
	}

	err = db.Commit()
	if err != nil {
		return "", types.ErrInternal
	}

	return hash, nil
}
