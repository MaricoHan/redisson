package service

import (
	"context"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"strconv"
)

type Nft struct {
	base *Base
}

func NewNft() *Nft {
	return &Nft{}
}
func (svc *Nft) EditNftByIndex(params dto.EditNftByIndexP) (string, error) {
	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return "", types.ErrMysqlConn
	}
	// get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, params.Index)).One(context.Background(), db)

	if err != nil {
		return "", types.ErrInternal
	}
	if tNft == nil || tNft.Status == "burned" {
		return "", types.ErrNftMissing
	}

	// create rawMsg
	msgEditNFT := nft.MsgEditNFT{
		Id:      strconv.FormatInt(int64(tNft.ID), 10),
		DenomId: tNft.ClassID,
		Name:    tNft.Name.String,
		URI:     params.Uri,
		Data:    params.Data,
		Sender:  params.Sender,
	}

	// lock the NFT
	tNft.Status = "pendding"
	_, err = tNft.Update(context.Background(), db, boil.Infer())
	if err != nil {
		return "", err
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)

	// Tx into database
	err = svc.base.TxIntoDataBase(params.AppID, txHash, signedData, "edit_nft", "undo")
	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, nil

}

func (svc *Nft) EditNftByBatch(params dto.EditNftByBatchP) (string, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return "", types.ErrMysqlConn
	}
	var msgEditNFTs sdktype.Msgs

	// create rawMsgs
	for _, EditNft := range params.EditNfts { // create every rawMsg
		// get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, EditNft.Index)).One(context.Background(), db)
		if err != nil {
			return "", err
		}
		msgEditNFT := nft.MsgEditNFT{
			Id:      strconv.FormatInt(int64(tNft.ID), 10),
			DenomId: tNft.ClassID,
			Name:    tNft.Name.String,
			URI:     EditNft.Uri,
			Data:    EditNft.Data,
			Sender:  params.Sender,
		}
		msgEditNFTs = append(msgEditNFTs, &msgEditNFT)

		// lock the NFT
		tNft.Status = "pendding"
		_, err = tNft.Update(context.Background(), db, boil.Infer())
		if err != nil {
			return "", err
		}
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(msgEditNFTs, baseTx)

	// Tx into database
	err = svc.base.TxIntoDataBase(params.AppID, txHash, signedData, "edit_nft_batch", "undo")
	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, nil
}

func (svc *Nft) DeleteNftByIndex(params dto.DeleteNftByIndexP) (string, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return "", types.ErrMysqlConn
	}

	// get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, params.Index)).One(context.Background(), db)
	if err != nil {
		return "", err
	}
	if tNft == nil || tNft.Status == "burned" {
		return "", types.ErrNftMissing
	}
	if tNft.Status == "pendding" {
		return "", types.ErrNftBurnPend
	}

	// lock the NFT
	tNft.Status = "pendding"
	_, err = tNft.Update(context.Background(), db, boil.Infer())
	if err != nil {
		return "", err
	}

	// create rawMsg
	msgBurnNFT := nft.MsgBurnNFT{
		Id:      tNft.NFTID,
		DenomId: tNft.ClassID,
		Sender:  params.Sender,
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgBurnNFT}, baseTx)

	// Tx into database
	err = svc.base.TxIntoDataBase(params.AppID, txHash, signedData, "burn_nft", "undo")
	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, nil
}
func (svc *Nft) DeleteNftByBatch(params dto.DeleteNftByBatchP) (string, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return "", types.ErrMysqlConn
	}

	var msgBurnNFTs sdktype.Msgs

	// create rawMsgs
	for _, index := range params.Indices { // create every rawMsg
		//get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, index)).One(context.Background(), db)
		if err != nil {
			return "", err
		}
		if tNft == nil || tNft.Status == "burned" {
			return "", types.ErrNftMissing
		}
		if tNft.Status == "pendding" {
			return "", types.ErrNftBurnPend
		}

		// create rawMsg
		msgBurnNFT := nft.MsgBurnNFT{
			Id:      tNft.NFTID,
			DenomId: tNft.ClassID,
			Sender:  params.Sender,
		}
		msgBurnNFTs = append(msgBurnNFTs, &msgBurnNFT)

		// lock the NFT
		tNft.Status = "pendding"
		_, err = tNft.Update(context.Background(), db, boil.Infer())
		if err != nil {
			return "", err
		}
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(msgBurnNFTs, baseTx)

	// Tx into database
	err = svc.base.TxIntoDataBase(params.AppID, txHash, signedData, "edit_nft_batch", "undo")
	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, nil
}

func (svc *Nft) NftByIndex(params dto.NftByIndexP) (*dto.NftByIndexP, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return nil, types.ErrMysqlConn
	}
	//get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, params.Index)).One(context.Background(), db)
	//get class by class_id
	class, err := models.TClasses(qm.Where("class_id=?", params.ClassId)).One(context.Background(), db)
	result := &dto.NftByIndexP{
		Id:          strconv.FormatInt(int64(tNft.ID), 10),
		Index:       tNft.Index,
		Name:        tNft.Name.String,
		ClassId:     tNft.ClassID,
		ClassName:   class.Name.String,
		ClassSymbol: class.Symbol.String,
		Uri:         tNft.URI.String,
		UriHash:     tNft.URIHash.String,
		Data:        tNft.Metadata.String,
		Owner:       tNft.Owner,
		Status:      tNft.Status,
		TxHash:      tNft.TXHash,
		TimeStamp:   tNft.Timestamp.Time.String(),
	}

	return result, nil
}
