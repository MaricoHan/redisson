package service

import (
	"context"
	"fmt"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
	"strconv"
	"strings"
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

func (svc *Nft) NftOperationHistoryByIndex(params dto.NftOperationHistoryByIndexP) (*dto.BNftOperationHistoryByIndexRes, error) {
	result := &dto.BNftOperationHistoryByIndexRes{}
	result.Offset = params.Offset
	result.Limit = params.Limit
	result.OperationRecords = []*dto.OperationRecord{}
	//nft, err := models.TNFTS(
	//	models.TNFTWhere.AppID.EQ(params.AppID),
	//	models.TNFTWhere.ClassID.EQ(params.ClassID),
	//	models.TNFTWhere.Index.EQ(params.Index),
	//	).OneG(context.Background())
	//if err != nil {
	//	return nil, types.ErrMysqlConn
	//}

	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TMSGS),
		qm.Select(models.TMSGColumns.TXHash,
			models.TMSGColumns.Operation,
			models.TMSGColumns.Signer,
			models.TMSGColumns.Recipient,
			models.TMSGColumns.Timestamp),
		models.TMSGWhere.AppID.EQ(params.AppID),
	}
	//if params.Txhash != "" {
	//	queryMod = append(queryMod, models.TMSGWhere.TXHash.EQ(params.Txhash))
	//}else {
	//	queryMod = append(queryMod, models.TMSGWhere.TXHash.EQ(nft.TXHash))
	//}
	////否则查询该nft的所有hash
	if params.Signer != "" {
		queryMod = append(queryMod, models.TMSGWhere.Signer.EQ(params.Signer))
	}
	if params.Operation != "" {
		queryMod = append(queryMod, models.TMSGWhere.Operation.EQ(params.Operation))
	}
	if params.StartDate != nil {
		queryMod = append(queryMod, models.TMSGWhere.CreateAt.GTE(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TMSGWhere.CreateAt.LTE(*params.EndDate))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s desc", models.TMSGWhere.CreateAt)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TMSGWhere.CreateAt)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}

	var modelResults []*models.TMSG
	total, err := modext.PageQuery(
		context.Background(),
		orm.GetDB(),
		queryMod,
		&modelResults,
		params.Offset,
		params.Limit,
	)
	if err != nil {
		// records not exist
		if strings.Contains(err.Error(), "records not exist") {
			return result, nil
		}

		return nil, types.ErrMysqlConn
	}

	result.TotalCount = total
	var operationRecords []*dto.OperationRecord
	for _, modelResult := range modelResults {
		var operationRecord = &dto.OperationRecord{
			Txhash:    modelResult.TXHash,
			Operation: modelResult.Operation,
			Signer:    modelResult.Signer,
			Recipient: modelResult.Recipient.String,
			Timestamp: modelResult.Timestamp.Time.String(),
		}
		operationRecords = append(operationRecords, operationRecord)
	}
	result.OperationRecords = operationRecords

	return result, nil
}
