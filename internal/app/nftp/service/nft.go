package service

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/friendsofgo/errors"

	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
)

const nftp = "nftp"

type Nft struct {
	base *Base
}

func NewNft(base *Base) *Nft {
	return &Nft{base: base}
}

func (svc *Nft) CreateNfts(params dto.CreateNftsRequest) ([]string, error) {
	var err error
	var txHash *string
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		classOne, err := models.TClasses(
			models.TClassWhere.AppID.EQ(params.AppID),
			models.TClassWhere.ClassID.EQ(params.ClassId),
		).One(context.Background(), exec)
		if err != nil {
			return types.ErrNftClassDetailsGet
		}
		if classOne.Status != models.TNFTSStatusActive {
			return types.ErrClassStatus
		}

		offSet := classOne.Offset
		var msgs sdktype.Msgs
		for i := 1; i <= params.Amount; i++ {
			index := int(offSet) + i
			nftId := nftp + strings.ToLower(hex.EncodeToString(tmhash.Sum([]byte(params.ClassId)))) + strconv.Itoa(index)
			if params.Recipient == "" {
				params.Recipient = classOne.Owner
			} else {
				acc, err := models.TAccounts(
					models.TAccountWhere.AppID.EQ(params.AppID),
					models.TAccountWhere.Address.EQ(params.Recipient)).OneG(context.Background())
				if err != nil {
					return types.ErrParams
				}
				if acc == nil {
					return types.ErrParams
				}
			}
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
		baseTx := svc.base.CreateBaseTx(classOne.Owner, "")
		originData, thash, err := svc.base.BuildAndSign(msgs, baseTx)
		if err != nil {
			log.Debug("create nfts", "buildandsign error:", err.Error())
			return err
		}
		//validate tx
		txone, err := svc.base.ValidateTx(thash)
		if err != nil {
			return err
		}
		if txone != nil && txone.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txone.ID)
			originData, thash, err = svc.base.BuildAndSign(msgs, baseTx)
			if err != nil {
				log.Debug("create nfts", "buildandsign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}
		txHash = &thash
		//transferInfo
		ttx := models.TTX{
			AppID:         params.AppID,
			Hash:          thash,
			Timestamp:     null.Time{Time: time.Now()},
			OriginData:    null.BytesFromPtr(&originData),
			OperationType: models.TTXSOperationTypeMintNFT,
			Status:        models.TTXSStatusUndo,
		}
		err = ttx.Insert(context.Background(), exec, boil.Infer())
		if err != nil {
			return types.ErrTxMsgInsert
		}
		classOne.Status = models.TTXSStatusPending
		classOne.Offset = classOne.Offset + uint64(params.Amount)
		classOne.LockedBy = null.Uint64FromPtr(&ttx.ID)
		ok, err := classOne.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			return types.ErrNftClassesSet
		}
		if ok != 1 {
			return types.ErrNftClassesSet
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return []string{*txHash}, nil
}

func (svc *Nft) EditNftByIndex(params dto.EditNftByIndexP) (string, error) {

	// get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(params.Index)).One(context.Background(), boil.GetContextDB())

	// internal error：500
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return "", types.ErrInternal
	}

	// nft does not exist ：404
	if tNft == nil || tNft.Status != models.TNFTSStatusActive {
		return "", types.ErrNftMissing
	}

	// judge whether the Caller is the owner：400
	if params.Sender != tNft.Owner {
		return "", types.ErrNotOwner
	}
	// judge whether the Caller is one of the APP's address：400
	if tNft.AppID != params.AppID {
		return "", types.ErrNoPermission
	}

	nftName := "[do-not-modify]"
	if params.Name != "" {
		nftName = params.Name
	}

	// create rawMsg
	msgEditNFT := nft.MsgEditNFT{
		Id:      tNft.NFTID,
		DenomId: tNft.ClassID,
		Name:    nftName,
		URI:     params.Uri,
		Data:    params.Data,
		Sender:  params.Sender,
		UriHash: "[do-not-modify]",
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)
	if err != nil {
		log.Debug("edit nft by index", "BuildAndSign error:", err.Error())
		return "", err
	}
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//validate tx
		txone, err := svc.base.ValidateTx(txHash)
		if err != nil {
			return err
		}
		if txone != nil && txone.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txone.ID)
			signedData, txHash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)
			if err != nil {
				log.Debug("edit nft by index", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		// Tx into database
		txId, err := svc.base.TxIntoDataBase(params.AppID, txHash, signedData, models.TTXSOperationTypeEditNFT, models.TTXSStatusUndo, exec)
		if err != nil {
			log.Debug("edit nft by index", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFT
		tNft.Status = models.TNFTSStatusPending
		tNft.LockedBy = null.Uint64From(txId)
		_, err = tNft.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, nil
}

func (svc *Nft) EditNftByBatch(params dto.EditNftByBatchP) (string, error) {

	// create rawMsgs
	var msgEditNFTs sdktype.Msgs
	for _, EditNft := range params.EditNfts { // create every rawMsg
		// get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(EditNft.Index)).One(context.Background(), boil.GetContextDB())
		// internal error：500
		if err != nil && errors.Cause(err) != sql.ErrNoRows {
			log.Error("edit nft by batch", "internal error:", err.Error())
			return "", types.ErrInternal
		}
		// nft does not exist or status is not active：400
		if tNft == nil {
			return "", types.ErrNftMissing
		}
		if tNft.Status != models.TNFTSStatusActive {
			return "", types.ErrNftStatus
		}

		// judge whether the Caller is the owner：400
		if params.Sender != tNft.Owner {
			return "", types.ErrNotOwner
		}
		// judge whether the Caller is one of the APP's address：400
		if tNft.AppID != params.AppID {
			return "", types.ErrNoPermission
		}

		nftName := "[do-not-modify]"
		if EditNft.Name != "" {
			nftName = EditNft.Name
		}
		msgEditNFT := nft.MsgEditNFT{
			Id:      tNft.NFTID,
			DenomId: tNft.ClassID,
			Name:    nftName,
			URI:     EditNft.Uri,
			Data:    EditNft.Data,
			Sender:  params.Sender,
			UriHash: "[do-not-modify]",
		}
		msgEditNFTs = append(msgEditNFTs, &msgEditNFT)
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(msgEditNFTs, baseTx)
	if err != nil {
		log.Debug("edit nft by batch", "BuildAndSign error:", err.Error())
		return "", err
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//validate tx
		txone, err := svc.base.ValidateTx(txHash)
		if err != nil {
			return err
		}
		if txone != nil && txone.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txone.ID)
			signedData, txHash, err = svc.base.BuildAndSign(msgEditNFTs, baseTx)
			if err != nil {
				log.Debug("edit nft by batch", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		// Tx into database
		txId, err := svc.base.TxIntoDataBase(params.AppID, txHash, signedData, models.TTXSOperationTypeEditNFTBatch, models.TTXSStatusUndo, exec)

		if err != nil {
			log.Debug("edit nft by batch", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFTs
		for _, EditNft := range params.EditNfts { // create every rawMsg
			tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(EditNft.Index)).One(context.Background(), exec)
			tNft.Status = models.TNFTSStatusPending
			tNft.LockedBy = null.Uint64From(txId)
			// update
			_, err = tNft.Update(context.Background(), exec, boil.Infer())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, nil
}

func (svc *Nft) DeleteNftByIndex(params dto.DeleteNftByIndexP) (string, error) {

	// get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(params.Index)).One(context.Background(), boil.GetContextDB())
	// internal error：500
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return "", types.ErrInternal
	}
	// nft does not exist or status is not active：404
	if tNft == nil || tNft.Status != models.TNFTSStatusActive {
		return "", types.ErrNftMissing
	}
	// pending：400
	if tNft.Status == models.TNFTSStatusPending {
		return "", types.ErrNftBurnPend
	}

	// judge whether the Caller is the owner：400
	if params.Sender != tNft.Owner {
		return "", types.ErrNotOwner
	}
	// judge whether the Caller is one of the APP's address：400
	if tNft.AppID != params.AppID {
		return "", types.ErrNoPermission
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
		log.Debug("delete nft by index", "BuildAndSign error:", err.Error())
		return "", err
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//validate tx
		txone, err := svc.base.ValidateTx(txHash)
		if err != nil {
			return err
		}
		if txone != nil && txone.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txone.ID)
			signedData, txHash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgBurnNFT}, baseTx)
			if err != nil {
				log.Debug("delete nft by index", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		// Tx into database
		txId, err := svc.base.TxIntoDataBase(params.AppID, txHash, signedData, models.TTXSOperationTypeBurnNFT, models.TTXSStatusUndo, exec)
		if err != nil {
			log.Debug("delete nft by index", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFT
		tNft.Status = models.TNFTSStatusPending
		tNft.LockedBy = null.Uint64From(txId)
		_, err = tNft.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, nil
}

func (svc *Nft) DeleteNftByBatch(params dto.DeleteNftByBatchP) (string, error) {
	// create rawMsgs
	var msgBurnNFTs sdktype.Msgs
	for _, index := range params.Indices { // create every rawMsg
		//get NFT by app_id,class_id and index
		tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(index)).One(context.Background(), boil.GetContextDB())
		// internal error：500
		if err != nil && errors.Cause(err) != sql.ErrNoRows {
			return "", types.ErrInternal
		}
		// nft does not exist or status is not active：400
		if tNft == nil || tNft.Status != models.TNFTSStatusActive {
			return "", types.ErrNftStatus
		}

		// judge whether the Caller is the owner
		if params.Sender != tNft.Owner {
			return "", types.ErrNotOwner
		}
		// judge whether the Caller is one of the APP's address
		if tNft.AppID != params.AppID {
			return "", types.ErrNoPermission
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
	if err != nil {
		log.Debug("delete nft by batch", "BuildAndSign error:", err.Error())
		return "", err
	}

	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//validate tx
		txone, err := svc.base.ValidateTx(txHash)
		if err != nil {
			return err
		}
		if txone != nil && txone.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txone.ID)
			signedData, txHash, err = svc.base.BuildAndSign(msgBurnNFTs, baseTx)
			if err != nil {
				log.Debug("delete nft by batch", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		// Tx into database
		txId, err := svc.base.TxIntoDataBase(params.AppID, txHash, signedData, models.TTXSOperationTypeBurnNFTBatch, models.TTXSStatusUndo, exec)
		if err != nil {
			log.Debug("delete nft by batch", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFTs
		for _, index := range params.Indices { // lock every nft
			tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(index)).One(context.Background(), exec)
			tNft.Status = models.TNFTSStatusPending
			tNft.LockedBy = null.Uint64From(txId)
			_, err = tNft.Update(context.Background(), exec, boil.Infer())
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	// return the txHash
	return txHash, nil
}

func (svc *Nft) NftByIndex(params dto.NftByIndexP) (*dto.NftR, error) {
	// get NFT by app_id,class_id and index
	tNft, err := models.TNFTS(models.TNFTWhere.AppID.EQ(params.AppID), models.TNFTWhere.ClassID.EQ(params.ClassId), models.TNFTWhere.Index.EQ(params.Index)).One(context.Background(), boil.GetContextDB())
	// internal error：500
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, types.ErrInternal
	}
	// nft does not exist
	if tNft == nil {
		return nil, types.ErrNftMissing
	}

	if !strings.Contains("active/burned", tNft.Status) {
		return nil, types.ErrNftStatus
	}

	// get class by class_id
	class, err := models.TClasses(models.TClassWhere.ClassID.EQ(params.ClassId)).One(context.Background(), boil.GetContextDB())
	// internal error：500
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, types.ErrInternal
	}
	if class == nil {
		return nil, types.ErrNftClassStatus
	}

	result := &dto.NftR{
		Id:          tNft.NFTID,
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
		Timestamp:   tNft.Timestamp.Time.String(),
	}

	return result, nil
}

func (svc *Nft) NftOperationHistoryByIndex(params dto.NftOperationHistoryByIndexP) (*dto.BNftOperationHistoryByIndexRes, error) {
	result := &dto.BNftOperationHistoryByIndexRes{
		PageRes: dto.PageRes{
			Offset:     params.Offset,
			Limit:      params.Limit,
			TotalCount: 0,
		},
		OperationRecords: nil,
	}
	res, err := models.TNFTS(
		models.TNFTWhere.AppID.EQ(params.AppID),
		models.TNFTWhere.ClassID.EQ(params.ClassID),
		models.TNFTWhere.Index.EQ(params.Index),
	).OneG(context.Background())
	if err != nil {
		return nil, types.ErrGetNftOperationDetails
	}

	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TMSGS),
		qm.Select(models.TMSGColumns.TXHash,
			models.TMSGColumns.Operation,
			models.TMSGColumns.Signer,
			models.TMSGColumns.Recipient,
			models.TMSGColumns.Timestamp),
		models.TMSGWhere.AppID.EQ(params.AppID),
		models.TMSGWhere.NFTID.EQ(null.StringFrom(res.NFTID)),
	}
	if params.Txhash != "" {
		queryMod = append(queryMod, models.TMSGWhere.TXHash.EQ(params.Txhash))
	}
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
			orderBy = fmt.Sprintf("%s DESC", models.TMSGColumns.CreateAt)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TMSGColumns.CreateAt)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}
	var modelResults []*models.TMSG
	total, err := modext.PageQueryByOffset(
		context.Background(),
		orm.GetDB(),
		queryMod,
		&modelResults,
		int(params.Offset),
		int(params.Limit),
	)
	if err != nil {
		// records not exist
		if strings.Contains(err.Error(), "records not exist") {
			return result, nil
		}
		return nil, types.ErrGetNftOperationDetails
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
	var err error
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
		queryMod = append(queryMod, models.TNFTWhere.Timestamp.GTE(null.TimeFromPtr(params.StartDate)))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TNFTWhere.Timestamp.LTE(null.TimeFromPtr(params.EndDate)))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "ID_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TNFTColumns.NFTID)
		case "ID_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TNFTColumns.NFTID)
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TNFTColumns.Timestamp)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TNFTColumns.Timestamp)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}

	var modelResults []*models.TNFT
	var total int64
	var classByIds []*dto.NftClassByIds
	classIds := []string{}
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		total, err = modext.PageQueryByOffset(
			context.Background(),
			exec,
			queryMod,
			&modelResults,
			int(params.Offset),
			int(params.Limit),
		)
		if err != nil {
			return err
		}

		tempMap := map[string]byte{} // 存放不重复主键
		for _, m := range modelResults {
			l := len(tempMap)
			tempMap[m.ClassID] = 0
			if len(tempMap) != l {
				classIds = append(classIds, m.ClassID)
			}
		}
		qMod := []qm.QueryMod{
			qm.From(models.TableNames.TClasses),
			qm.Select(models.TClassColumns.ClassID, models.TClassColumns.Name, models.TClassColumns.Symbol),
			models.TClassWhere.ClassID.IN(classIds),
		}

		err = models.NewQuery(qMod...).Bind(context.Background(), exec, &classByIds)
		return err
	})

	if err != nil {
		// records not exist
		if strings.Contains(err.Error(), "records not exist") {
			return result, nil
		}
		return nil, types.ErrMysqlConn
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
			Timestamp: modelResult.Timestamp.Time.String(),
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
