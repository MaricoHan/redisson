package service

import (
	"context"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/volatiletech/null/v8"
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
	tNft.Status = "undo"
	_, err = tNft.Update(context.Background(), db, boil.Infer())
	if err != nil {
		return "", err
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)

	// Tx into database
	ttx := models.TTX{
		AppID:         params.AppID,
		Hash:          txHash,
		OriginData:    null.BytesFrom(signedData),
		OperationType: "edit_nft",
		Status:        "undo",
	}
	err = ttx.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, err

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
		tNft.Status = "undo"
		_, err = tNft.Update(context.Background(), db, boil.Infer())
		if err != nil {
			return "", err
		}
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(msgEditNFTs, baseTx)

	// Tx into database
	ttx := models.TTX{
		AppID:         params.AppID,
		Hash:          txHash,
		OriginData:    null.BytesFrom(signedData),
		OperationType: "edit_nft_batch",
		Status:        "undo",
	}
	err = ttx.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, err
}

func (svc *Nft) DeleteNftByIndex(params dto.DeleteNftByIndexP) (string, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return "", types.ErrMysqlConn
	}

	//get NFT by app_id,class_id and index
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
	// create rawMsg
	msgBurnNFT := nft.MsgBurnNFT{
		Id:      tNft.NFTID,
		DenomId: tNft.ClassID,
		Sender:  params.Owner,
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Owner, "")
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgBurnNFT}, baseTx)

	// Tx into database
	ttx := models.TTX{
		AppID:         params.AppID,
		Hash:          txHash,
		OriginData:    null.BytesFrom(signedData),
		OperationType: "burn_nft",
		Status:        "undo",
	}

	err = ttx.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		return "", err
	}

	// lock the NFT
	tNft.Status = "undo"
	_, err = tNft.Update(context.Background(), db, boil.Infer())
	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, err
}
func (svc *Nft) DeleteNftByBatch(params dto.DeleteNftByBatchP) (int64, error) {

	// get database object
	db, err := orm.GetDB().Begin()
	if err != nil {
		return 0, types.ErrMysqlConn
	}
	var rowsAff int64
	for _, index := range params.Indices { //burn every NFT
		//get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(qm.Where("app_id=? AND class_id=? AND index=?", params.AppID, params.ClassId, index)).One(context.Background(), db)
		if err != nil {
			return rowsAff, err
		}
		if tNft == nil || tNft.Status == "burned" {
			return rowsAff, types.ErrNftMissing
		}
		if tNft.Status == "pendding" {
			return rowsAff, types.ErrNftBurnPend
		}

		//just burn
		tNft.Status = "burned"
		i, err := tNft.Update(context.Background(), db, boil.Infer())
		rowsAff += i

		if err != nil {
			return rowsAff, err
		}
	}

	//return the affected rows amount
	return rowsAff, nil
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

func (svc *Nft) NftOperationHistoryByIndex(params dto.NftOperationHistoryByIndexP) *dto.NftOperationHistoryByIndexRes {
	return nil
}
