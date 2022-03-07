package service

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
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

func (svc *Nft) CreateNfts(params dto.CreateNftsP) (*dto.TxRes, error) {
	var err error
	var taskId string
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		classOne, err := models.TClasses(
			models.TClassWhere.ProjectID.EQ(params.ProjectID),
			models.TClassWhere.ClassID.EQ(params.ClassId),
		).One(context.Background(), exec)

		if err != nil {
			if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), SqlNotFound) {
				//404
				return types.ErrNotFound
			}
			//500
			log.Error("create nfts", "query class error:", err.Error())
			return types.ErrInternal
		}

		//400
		if classOne.Status != models.TNFTSStatusActive {
			return types.ErrNftClassStatus
		}

		// ValidateSigner
		if err := svc.base.ValidateSigner(classOne.Owner, params.ProjectID); err != nil {
			return err
		}

		offSet := classOne.Offset
		var msgs sdktype.Msgs
		for i := 1; i <= params.Amount; i++ {
			index := int(offSet) + i
			nftId := nftp + strings.ToLower(hex.EncodeToString(tmhash.Sum([]byte(params.ClassId)))) + strconv.Itoa(index)
			if params.Recipient == "" {
				//默认为 NFT 类别的权属者地址
				params.Recipient = classOne.Owner
			}
			// ValidateRecipient
			if err := svc.base.ValidateRecipient(params.Recipient, params.ProjectID); err != nil {
				return err
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
		originData, tHash, _ := svc.base.BuildAndSign(msgs, baseTx)
		baseTx.Gas = svc.base.mintNftsGas(originData, uint64(params.Amount))
		err = svc.base.GasThan(classOne.Owner, params.ChainID, baseTx.Gas, params.PlatFormID)
		if err != nil {
			return types.NewAppError(types.RootCodeSpace, types.ErrGasNotEnough, err.Error())
		}
		originData, tHash, err = svc.base.BuildAndSign(msgs, baseTx)
		if err != nil {
			log.Debug("create nfts", "buildandsign error:", err.Error())
			return types.ErrBuildAndSign
		}

		//validate tx
		txOne, err := svc.base.ValidateTx(tHash)
		if err != nil {
			return err
		}
		if txOne != nil && txOne.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txOne.ID)
			originData, tHash, err = svc.base.BuildAndSign(msgs, baseTx)
			if err != nil {
				log.Debug("create nfts", "buildandsign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		msgsBytes, _ := json.Marshal(msgs)
		code := fmt.Sprintf("%s%s%s", classOne.Owner, models.TTXSOperationTypeMintNFT, time.Now().String())
		taskId = svc.base.EncodeData(code)
		ttx := models.TTX{
			ProjectID:     params.ProjectID,
			Hash:          tHash,
			Timestamp:     null.Time{Time: time.Now()},
			Message:       null.JSONFrom(msgsBytes),
			Sender:        null.StringFrom(classOne.Owner),
			TaskID:        null.StringFrom(taskId),
			GasUsed:       null.Int64From(int64(baseTx.Gas)),
			OriginData:    null.BytesFromPtr(&originData),
			OperationType: models.TTXSOperationTypeMintNFT,
			Status:        models.TTXSStatusUndo,
			Tag:           null.JSONFrom(params.Tag),
		}

		err = ttx.Insert(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Error("create nft", "ttx insert error: ", err)
			return types.ErrInternal
		}

		//class locked
		classOne.Status = models.TTXSStatusPending
		classOne.Offset = classOne.Offset + uint64(params.Amount)
		classOne.LockedBy = null.Uint64FromPtr(&ttx.ID)
		ok, err := classOne.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Error("create nft", "class status update error: ", err)
			return types.ErrInternal
		}
		if ok != 1 {
			return types.ErrInternal
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return &dto.TxRes{TaskId: taskId}, nil
}

func (svc *Nft) EditNftByNftId(params dto.EditNftByNftIdP) (*dto.TxRes, error) {
	// ValidateSigner
	if err := svc.base.ValidateSigner(params.Sender, params.ProjectID); err != nil {
		return nil, err
	}
	tNft, err := models.TNFTS(models.TNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TNFTWhere.ClassID.EQ(params.ClassId),
		models.TNFTWhere.NFTID.EQ(params.NftId),
		models.TNFTWhere.Owner.EQ(params.Sender)).
		One(context.Background(), boil.GetContextDB())

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), SqlNotFound) {
			//404
			return nil, types.ErrNotFound
		}
		//500
		log.Error("edit nft by nftId", "query nft error:", err.Error())
		return nil, types.ErrInternal
	}
	//404
	if tNft.Status == models.TNFTSStatusBurned {
		return nil, types.ErrNotFound
	}
	//400
	if tNft.Status != models.TNFTSStatusActive {
		return nil, types.ErrNftStatus
	}

	//非必填保留数据
	uri := params.Uri
	if uri == "" {
		uri = tNft.URI.String
	}
	data := params.Data
	if data == "" {
		data = tNft.Metadata.String
	}

	// create rawMsg
	msgEditNFT := nft.MsgEditNFT{
		Id:      tNft.NFTID,
		DenomId: tNft.ClassID,
		Name:    params.Name,
		URI:     uri,
		Data:    data,
		Sender:  params.Sender,
		UriHash: "[do-not-modify]",
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)

	// get gas
	nftLen := svc.base.lenOfNft(tNft)
	baseTx.Gas = svc.base.editNftGas(nftLen, uint64(len(signedData)))
	err = svc.base.GasThan(params.Sender, params.ChainID, baseTx.Gas, params.PlatFormID)
	if err != nil {
		return nil, types.NewAppError(types.RootCodeSpace, types.ErrGasNotEnough, err.Error())
	}
	signedData, txHash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)

	if err != nil {
		log.Debug("edit nft by nftId", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		//validate tx
		txOne, err := svc.base.ValidateTx(txHash)
		if err != nil {
			return err
		}
		if txOne != nil && txOne.Status == models.TTXSStatusFailed {
			baseTx.Memo = fmt.Sprintf("%d", txOne.ID)
			signedData, txHash, err = svc.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)
			if err != nil {
				log.Debug("edit nft by nftId", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		// Tx into database
		messageByte, _ := json.Marshal(msgEditNFT)
		code := fmt.Sprintf("%s%s%s", params.Sender, models.TTXSOperationTypeEditNFT, time.Now().String())
		taskId = svc.base.EncodeData(code)

		// Tx into database
		txId, err := svc.base.UndoTxIntoDataBase(params.Sender, models.TTXSOperationTypeEditNFT, taskId, txHash,
			params.ProjectID, signedData, messageByte, params.Tag, int64(baseTx.Gas), exec)
		if err != nil {
			log.Debug("edit nft by nftId", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFT
		tNft.Status = models.TNFTSStatusPending
		tNft.LockedBy = null.Uint64From(txId)
		ok, err := tNft.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			return types.ErrInternal
		}
		if ok != 1 {
			return types.ErrInternal
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &dto.TxRes{TaskId: taskId}, nil
}

func (svc *Nft) EditNftByBatch(params dto.EditNftByBatchP) (*dto.TxRes, error) {
	// create rawMsgs
	var msgEditNFTs sdktype.Msgs
	var nftsLen uint64
	for i, EditNft := range params.EditNfts {
		tNft, err := models.TNFTS(models.TNFTWhere.ProjectID.EQ(params.ProjectID),
			models.TNFTWhere.ClassID.EQ(params.ClassId),
			models.TNFTWhere.Owner.EQ(params.Sender),
			models.TNFTWhere.NFTID.EQ(EditNft.NftId)).
			One(context.Background(), boil.GetContextDB())

		if err != nil {
			if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), SqlNotFound) {
				//404
				return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftFound)
			}
			//500
			log.Error("edit nft by batch", "query nft error:", err.Error())
			return nil, types.ErrInternal
		}

		if tNft.Status == models.TNFTSStatusBurned {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftFound)
		}

		//400
		if tNft.Status != models.TNFTSStatusActive {
			return nil, types.NewAppError(types.RootCodeSpace, types.NftStatusAbnormal, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftStatusMsg)
		}
		// get nftLen
		nftLen := svc.base.lenOfNft(tNft)
		nftsLen += nftLen

		//非必填保留数据
		uri := EditNft.Uri
		if uri == "" {
			uri = tNft.URI.String
		}
		data := EditNft.Data
		if data == "" {
			data = tNft.Metadata.String
		}

		msgEditNFT := nft.MsgEditNFT{
			Id:      tNft.NFTID,
			DenomId: tNft.ClassID,
			Name:    EditNft.Name,
			URI:     uri,
			Data:    data,
			Sender:  params.Sender,
			UriHash: "[do-not-modify]",
		}
		msgEditNFTs = append(msgEditNFTs, &msgEditNFT)
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := svc.base.BuildAndSign(msgEditNFTs, baseTx)

	// set gas
	baseTx.Gas = svc.base.editBatchNftGas(nftsLen, uint64(len(signedData)))
	signedData, txHash, err = svc.base.BuildAndSign(msgEditNFTs, baseTx)

	if err != nil {
		log.Debug("edit nft by batch", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
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
		messageByte, _ := json.Marshal(msgEditNFTs)
		code := fmt.Sprintf("%s%s%s", params.Sender, models.TTXSOperationTypeEditNFTBatch, time.Now().String())
		taskId = svc.base.EncodeData(code)
		txId, err := svc.base.UndoTxIntoDataBase(params.Sender, models.TTXSOperationTypeEditNFTBatch, taskId, txHash,
			params.ProjectID, signedData, messageByte, nil, int64(baseTx.Gas), exec)
		if err != nil {
			log.Debug("edit nft by batch", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFTs
		for _, EditNft := range params.EditNfts { // create every rawMsg
			tNft, err := models.TNFTS(models.TNFTWhere.ProjectID.EQ(params.ProjectID),
				models.TNFTWhere.ClassID.EQ(params.ClassId),
				models.TNFTWhere.NFTID.EQ(EditNft.NftId)).
				One(context.Background(), exec)
			tNft.Status = models.TNFTSStatusPending
			tNft.LockedBy = null.Uint64From(txId)
			// update
			ok, err := tNft.Update(context.Background(), exec, boil.Infer())
			if err != nil {
				return types.ErrInternal
			}
			if ok != 1 {
				return types.ErrInternal
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.TxRes{TaskId: taskId}, nil
}

func (svc *Nft) DeleteNftByNftId(params dto.DeleteNftByNftIdP) (*dto.TxRes, error) {
	// ValidateSigner
	if err := svc.base.ValidateSigner(params.Sender, params.ProjectID); err != nil {
		return nil, err
	}

	tNft, err := models.TNFTS(models.TNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TNFTWhere.ClassID.EQ(params.ClassId),
		models.TNFTWhere.NFTID.EQ(params.NftId),
		models.TNFTWhere.Owner.EQ(params.Sender)).
		One(context.Background(), boil.GetContextDB())
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), SqlNotFound) {
			//404
			return nil, types.ErrNotFound
		}
		//500
		log.Error("delete nft by nftId", "query nft error:", err.Error())
		return nil, types.ErrInternal
	}

	// 404
	if tNft.Status == models.TNFTSStatusBurned {
		return nil, types.ErrNotFound
	}

	//400
	if tNft.Status != models.TNFTSStatusActive {
		return nil, types.ErrNftStatus
	}

	// create rawMsg
	msgBurnNFT := nft.MsgBurnNFT{
		Id:      tNft.NFTID,
		DenomId: tNft.ClassID,
		Sender:  params.Sender,
	}

	// build and sign transaction
	baseTx := svc.base.CreateBaseTx(params.Sender, "")

	nftLen := svc.base.lenOfNft(tNft)
	// set gas
	baseTx.Gas = svc.base.deleteNftGas(nftLen)
	err = svc.base.GasThan(params.Sender, params.ChainID, baseTx.Gas, params.PlatFormID)
	if err != nil {
		return nil, types.NewAppError(types.RootCodeSpace, types.ErrGasNotEnough, err.Error())
	}
	signedData, txHash, err := svc.base.BuildAndSign(sdktype.Msgs{&msgBurnNFT}, baseTx)

	if err != nil {
		log.Debug("delete nft by nftId", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
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
				log.Debug("delete nft by nftId", "BuildAndSign error:", err.Error())
				return types.ErrBuildAndSign
			}
		}

		// Tx into database
		messageByte, _ := json.Marshal(msgBurnNFT)
		code := fmt.Sprintf("%s%s%s", params.Sender, models.TTXSOperationTypeBurnNFT, time.Now().String())
		taskId = svc.base.EncodeData(code)
		// Tx into database
		txId, err := svc.base.UndoTxIntoDataBase(params.Sender, models.TTXSOperationTypeBurnNFT, taskId, txHash,
			params.ProjectID, signedData, messageByte, params.Tag, int64(baseTx.Gas), exec)

		if err != nil {
			log.Debug("delete nft by nftId", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFT
		tNft.Status = models.TNFTSStatusPending
		tNft.LockedBy = null.Uint64From(txId)
		ok, err := tNft.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			return types.ErrInternal
		}
		if ok != 1 {
			return types.ErrInternal
		}
		return err
	})

	if err != nil {
		return nil, err
	}
	return &dto.TxRes{TaskId: taskId}, nil
}

func (svc *Nft) DeleteNftByBatch(params dto.DeleteNftByBatchP) (*dto.TxRes, error) {
	// ValidateSigner
	if err := svc.base.ValidateSigner(params.Sender, params.ProjectID); err != nil {
		return nil, err
	}
	// create rawMsgs
	var msgBurnNFTs sdktype.Msgs
	var nftsLen uint64
	for i, nftId := range params.NftIds {
		tNft, err := models.TNFTS(models.TNFTWhere.ProjectID.EQ(params.ProjectID),
			models.TNFTWhere.ClassID.EQ(params.ClassId),
			models.TNFTWhere.NFTID.EQ(nftId),
			models.TNFTWhere.Owner.EQ(params.Sender)).
			One(context.Background(), boil.GetContextDB())
		if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
			(err != nil && strings.Contains(err.Error(), SqlNotFound)) {
			//400
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftFound)
		} else if err != nil {
			//500
			log.Error("delete nft by batch", "query nft error:", err.Error())
			return nil, types.ErrInternal
		}

		if tNft.Status == models.TNFTSStatusBurned {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftFound)
		}

		//400
		if tNft.Status != models.TNFTSStatusActive {
			return nil, types.NewAppError(types.RootCodeSpace, types.NftStatusAbnormal, "the "+fmt.Sprintf("%d", i+1)+"th "+types.ErrNftStatusMsg)
		}

		nftLen := svc.base.lenOfNft(tNft)
		nftsLen += nftLen

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
	// set gas
	baseTx.Gas = svc.base.deleteBatchNftGas(nftsLen, uint64(len(params.NftIds)))
	signedData, txHash, err = svc.base.BuildAndSign(msgBurnNFTs, baseTx)

	if err != nil {
		log.Debug("delete nft by batch", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
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
		messageByte, _ := json.Marshal(msgBurnNFTs)
		code := fmt.Sprintf("%s%s%s", params.Sender, models.TTXSOperationTypeBurnNFTBatch, time.Now().String())
		taskId = svc.base.EncodeData(code)
		// Tx into database
		txId, err := svc.base.UndoTxIntoDataBase(params.Sender, models.TTXSOperationTypeBurnNFTBatch, taskId, txHash,
			params.ProjectID, signedData, messageByte, nil, int64(baseTx.Gas), exec)
		if err != nil {
			log.Debug("delete nft by batch", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFTs
		for _, nftId := range params.NftIds { // lock every nft
			tNft, err := models.TNFTS(
				models.TNFTWhere.ProjectID.EQ(params.ProjectID),
				models.TNFTWhere.ClassID.EQ(params.ClassId),
				models.TNFTWhere.NFTID.EQ(nftId)).
				One(context.Background(), exec)
			tNft.Status = models.TNFTSStatusPending
			tNft.LockedBy = null.Uint64From(txId)
			ok, err := tNft.Update(context.Background(), exec, boil.Infer())
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
		return nil, err
	}

	return &dto.TxRes{TaskId: txHash}, nil
}

func (svc *Nft) NftByNftId(params dto.NftByNftIdP) (*dto.NftR, error) {
	// get NFT by app_id,class_id and nftId
	tNft, err := models.TNFTS(models.TNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TNFTWhere.ClassID.EQ(params.ClassId),
		models.TNFTWhere.NFTID.EQ(params.NftId)).
		One(context.Background(), boil.GetContextDB())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNotFound)) {
		//404
		return nil, types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("nft by nftId", "query nft error:", err.Error())
		return nil, types.ErrInternal
	}

	if strings.Contains(models.TNFTSStatusPending, tNft.Status) {
		return nil, types.ErrNftStatus
	}

	// get class by class_id
	class, err := models.TClasses(models.TClassWhere.ClassID.EQ(params.ClassId)).
		One(context.Background(), boil.GetContextDB())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNotFound)) {
		//404
		return nil, types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("nft by nftId", "query nft class error:", err.Error())
		return nil, types.ErrInternal
	}

	result := &dto.NftR{
		Id:          tNft.NFTID,
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

func (svc *Nft) NftOperationHistoryByNftId(params dto.NftOperationHistoryByNftIdP) (*dto.BNftOperationHistoryByNftIdRes, error) {
	result := &dto.BNftOperationHistoryByNftIdRes{
		PageRes: dto.PageRes{
			Offset:     params.Offset,
			Limit:      params.Limit,
			TotalCount: 0,
		},
		OperationRecords: []*dto.OperationRecord{},
	}
	res, err := models.TNFTS(
		models.TNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TNFTWhere.ClassID.EQ(params.ClassID),
		models.TNFTWhere.NFTID.EQ(params.NftId),
	).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNotFound)) {
		//404
		return nil, types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("query nft operation history", "query nft error:", err.Error())
		return nil, types.ErrInternal
	}

	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TMSGS),
		qm.Select(models.TMSGColumns.TXHash,
			models.TMSGColumns.Operation,
			models.TMSGColumns.Signer,
			models.TMSGColumns.Recipient,
			models.TMSGColumns.Timestamp),
		models.TMSGWhere.ProjectID.EQ(params.ProjectID),
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
		if strings.Contains(err.Error(), SqlNotFound) {
			return result, nil
		}
		return nil, types.ErrInternal
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
		models.TNFTWhere.ProjectID.EQ(params.ProjectID),
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
	var classIds []string
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
		if strings.Contains(err.Error(), SqlNotFound) {
			return result, nil
		}
		return nil, types.ErrInternal
	}

	result.TotalCount = total
	var nfts []*dto.Nft
	for _, modelResult := range modelResults {
		nft := &dto.Nft{
			Id:        modelResult.NFTID,
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
