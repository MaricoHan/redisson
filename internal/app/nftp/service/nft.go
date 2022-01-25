package service

import (
	"context"
	"encoding/hex"
	"fmt"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
	"time"
)

const nftp = "nftp"

type Nft struct {
	base *Base
}

func NewNft(base *Base) *Nft {
	return &Nft{base: base}
}

func (svc *Nft) CreateNfts(params dto.CreateNftsRequest) ([]string, error) {
	db, err := orm.GetDB().Begin()
	//platform address
	classOne, err := models.TClasses(
		models.TClassWhere.ClassID.EQ(params.ClassId),
	).One(context.Background(), db)
	if err != nil {
		db.Rollback()
		return nil, types.ErrGetAccountDetails
	}
	offSet := classOne.Offset
	//nftID := nftp + sha256(nftClassID)+index
	var msgs sdktype.Msgs
	for i := 1; i <= params.Amount; i++ {
		index := int(offSet) + i
		nftId := nftp + strings.ToUpper(hex.EncodeToString(tmhash.Sum([]byte(params.ClassId)))) + string(index)
		createNft := nft.MsgMintNFT{
			Id:        nftId,
			DenomId:   params.ClassId,
			Name:      params.Name,
			URI:       params.Uri,
			UriHash:   params.UriHash,
			Data:      params.Data,
			Sender:    classOne.Owner,
			Recipient: params.Recipient,
		}
		msgs = append(msgs, &createNft)
	}
	baseTx := svc.base.CreateBaseTx(classOne.Owner, defultKeyPassword)
	originData, txHash, err := svc.base.BuildAndSign(msgs, baseTx)
	if err != nil {
		return nil, err
	}

	//modify t_class offset
	tClass, err := models.TClasses(models.TClassWhere.AppID.EQ(params.AppID), models.TClassWhere.ClassID.EQ(params.ClassId)).One(context.Background(), db)
	if err != nil {
		db.Rollback()
		return nil, types.ErrNftClassesGet
	}
	tClass.Status = models.TTXSStatusPending
	tClass.Offset = tClass.Offset + uint64(params.Amount)

	//transferInfo
	ttx := models.TTX{
		AppID:         params.AppID,
		Hash:          txHash,
		Timestamp:     null.Time{Time: time.Now()},
		OriginData:    null.BytesFromPtr(&originData),
		OperationType: models.TTXSOperationTypeMintNFT,
		Status:        models.TTXSStatusUndo,
	}

	err = ttx.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		db.Rollback()
		return nil, types.ErrTxMsgInsert
	}

	tx, err := models.TTXS(qm.Where("hash=?", txHash)).One(context.Background(), boil.GetContextDB())
	if err != nil {
		db.Rollback()
		return nil, types.ErrTxMsgGet
	}

	tClass.LockedBy = tx.ID
	_, err = tClass.Update(context.Background(), db, boil.Infer())
	if err != nil {
		db.Rollback()
		return nil, types.ErrNftClassesSet
	}
	err = db.Commit()
	if err != nil {
		return nil, types.ErrInternal
	}

	var hashs []string
	hashs = append(hashs, txHash)
	return hashs, nil
}

func (svc *Nft) EditNftByIndex(params dto.EditNftByIndexP) (string, error) {

	// get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(params.Index)).One(context.Background(), boil.GetContextDB())
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

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)
	if err != nil {
		return "", err
	}

	// Tx into database
	txId, err := svc.base.TxIntoDataBase(params.AppID, txHash, signedData, "edit_nft", models.TTXSStatusUndo)
	if err != nil {
		return "", err
	}

	// lock the NFT
	tNft.Status = "pendding"
	tNft.LockedBy = null.Uint64From(txId)
	_, err = tNft.UpdateG(context.Background(), boil.Infer())
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

	// create rawMsgs
	var msgEditNFTs sdktype.Msgs
	for _, EditNft := range params.EditNfts { // create every rawMsg
		// get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(EditNft.Index)).One(context.Background(), boil.GetContextDB())

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

	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(msgEditNFTs, baseTx)
	if err != nil {
		return "", err
	}

	// Tx into database
	txId, err := svc.base.TxIntoDataBase(params.AppID, txHash, signedData, "edit_nft_batch", models.TTXSStatusUndo)
	if err != nil {
		return "", err
	}

	// lock the NFTs
	for _, EditNft := range params.EditNfts { // create every rawMsg
		tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(EditNft.Index)).One(context.Background(), boil.GetContextDB())
		tNft.Status = "pendding"
		tNft.LockedBy = null.Uint64From(txId)
		// update
		_, err = tNft.Update(context.Background(), db, boil.Infer())
		if err != nil {
			return "", err
		}
	}

	// commit
	err = db.Commit()
	if err != nil {
		return "", types.ErrInternal
	}

	// return the txHash
	return txHash, nil
}

func (svc *Nft) DeleteNftByIndex(params dto.DeleteNftByIndexP) (string, error) {

	// get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(params.Index)).One(context.Background(), boil.GetContextDB())
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

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgBurnNFT}, baseTx)
	if err != nil {
		return "", err
	}

	// Tx into database
	txId, err := svc.base.TxIntoDataBase(params.AppID, txHash, signedData, "burn_nft", models.TTXSStatusUndo)
	if err != nil {
		return "", err
	}

	// lock the NFT
	tNft.Status = "pendding"
	tNft.LockedBy = null.Uint64From(txId)
	_, err = tNft.UpdateG(context.Background(), boil.Infer())
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

	// create rawMsgs
	var msgBurnNFTs sdktype.Msgs
	for _, index := range params.Indices { // create every rawMsg
		//get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(index)).One(context.Background(), boil.GetContextDB())
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
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(msgBurnNFTs, baseTx)

	// Tx into database
	txId, err := svc.base.TxIntoDataBase(params.AppID, txHash, signedData, "burn_nft_batch", models.TTXSStatusUndo)
	if err != nil {
		return "", err
	}

	// lock the NFT
	for _, index := range params.Indices { // create every rawMsg
		tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(index)).One(context.Background(), boil.GetContextDB())
		tNft.Status = "pendding"
		tNft.LockedBy = null.Uint64From(txId)
		_, err = tNft.Update(context.Background(), db, boil.Infer())
		if err != nil {
			return "", err
		}
	}
	// commit
	err = db.Commit()
	if err != nil {
		return "", types.ErrInternal
	}

	// return the txHash
	return txHash, nil
}

func (svc *Nft) NftByIndex(params dto.NftByIndexP) (*dto.NftByIndexP, error) {

	// get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(params.Index)).One(context.Background(), boil.GetContextDB())
	if err != nil {
		return nil, err
	}

	// get class by class_id
	class, err := models.TClasses(models.TClassWhere.ClassID.EQ(params.ClassId)).One(context.Background(), boil.GetContextDB())

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

func (svc *Nft) Nfts(params dto.NftsP) (*dto.NftsRes, error) {
	db, err := orm.GetDB().Begin()
	result := &dto.NftsRes{}
	result.Offset = params.Offset
	result.Limit = params.Limit
	result.Nfts = []*dto.Nft{}
	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TNFTS),
		models.TNFTWhere.AppID.EQ(params.AppID),
	}
	if params.Id != "" {
		queryMod = append(queryMod, models.TNFTWhere.NFTID.EQ(params.Id))
	}
	if params.ClassId != "" {
		queryMod = append(queryMod, models.TNFTWhere.ClassID.EQ(params.ClassId))
	}
	if params.Owner != "" {
		queryMod = append(queryMod, models.TNFTWhere.Owner.EQ(params.Owner))
	}
	if params.TxHash != "" {
		queryMod = append(queryMod, models.TNFTWhere.TXHash.EQ(params.TxHash))
	}
	if params.Status != "" {
		queryMod = append(queryMod, models.TNFTWhere.Status.EQ(params.Status))
	}

	if params.StartDate != nil {
		queryMod = append(queryMod, models.TNFTWhere.CreateAt.GTE(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TNFTWhere.CreateAt.LTE(*params.EndDate))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s desc", models.TNFTColumns.CreateAt)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TNFTColumns.CreateAt)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}

	var modelResults []*models.TNFT
	total, err := modext.PageQueryByOffset(
		context.Background(),
		db,
		queryMod,
		&modelResults,
		int(params.Offset),
		int(params.Limit),
	)
	if err != nil {
		db.Rollback()
		// records not exist
		if strings.Contains(err.Error(), "records not exist") {
			return result, nil
		}
		return nil, types.ErrMysqlConn
	}

	classIds := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, m := range modelResults {
		l := len(tempMap)
		tempMap[m.ClassID] = 0
		if len(tempMap) != l {
			classIds = append(classIds, m.ClassID) //当元素不重复时，将元素添加到切片result中
		}
	}
	q1 := []qm.QueryMod{
		qm.From(models.TableNames.TClasses),
		qm.Select(models.TClassColumns.ClassID, models.TClassColumns.Name, models.TClassColumns.Symbol),
	}
	q1 = append(q1, models.TClassWhere.ClassID.IN(classIds))
	var classByIds []*dto.NftClassByIds
	err = models.NewQuery(q1...).Bind(context.Background(), db, &classByIds)
	if err != nil {
		db.Rollback()
		return nil, types.ErrInternal
	}
	err = db.Commit()
	if err != nil {
		return nil, types.ErrInternal
	}

	result.TotalCount = total
	var nfts []*dto.Nft
	for _, modelResult := range modelResults {
		nft := &dto.Nft{
			Id:        modelResult.NFTID,
			Index:     modelResult.Index,
			Name:      modelResult.Name.String,
			ClassId:   modelResult.ClassID,
			Uri:       modelResult.URI.String,
			Owner:     modelResult.Owner,
			Status:    modelResult.Status,
			TxHash:    modelResult.TXHash,
			TimeStamp: modelResult.Timestamp.Time.String(),
		}
		for _, class := range classByIds {
			if class.ClassId == modelResult.ClassID {
				nft.ClassName = class.Name
				nft.ClassSymbol = class.Symbol
			}
		}
		nfts = append(nfts, nft)
	}
	result.Nfts = nfts
	return result, nil
}
