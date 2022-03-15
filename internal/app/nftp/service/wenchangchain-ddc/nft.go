package wenchangchain_ddc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/friendsofgo/errors"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
	"strings"
	"time"
)

const nftp = "nftp"

type Nft struct {
	base *service.Base
}

func (n Nft) List(params dto.NftsP) (*dto.NftsRes, error) {
	panic("implement me")
}

func (n Nft) Create(params dto.CreateNftsP) (*dto.TxRes, error) {
	DDC721Service := service.DDCClient.GetDDC721Service(true, true)
	var taskId string
	err := modext.Transaction(func(exec boil.ContextExecutor) error {
		// query class
		class, err := models.TDDCClasses(
			models.TDDCClassWhere.ProjectID.EQ(params.ProjectID),
			models.TDDCClassWhere.ClassID.EQ(params.ClassId),
		).One(context.Background(), exec)
		if err != nil {
			if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), service.SqlNotFound) {
				//404
				return types.ErrNotFound
			}
			//500
			log.Error("create ddc", "query class error:", err.Error())
			return types.ErrInternal
		}

		//400
		if class.Status != models.TDDCClassesStatusActive {
			return types.ErrNftClassStatus
		}

		// ValidateSigner
		if err := n.base.ValidateSigner(class.Owner, params.ProjectID); err != nil {
			return err
		}

		if params.Recipient == "" {
			//默认为 NFT 类别的权属者地址
			params.Recipient = class.Owner
		}

		// ValidateRecipient
		if err := n.base.ValidateRecipient(params.Recipient, params.ProjectID); err != nil {
			return err
		}

		//taskId
		code := fmt.Sprintf("%s%s%s", class.Owner, models.TTXSOperationTypeMintNFT, time.Now().String())
		taskId = n.base.EncodeData(code)

		//签名后的交易计算动态 gas
		res, err := DDC721Service.SafeMint(bind.TransactOpts{}, class.Owner, params.Recipient, params.Uri, []byte(params.Data))
		//originData, err := DDC721Service.SignTx(common.HexToAddress(class.Owner), res.RawTx)

		//离线数据组装为 msg 入表传到 block sync
		createNft := nft.MsgMintNFT{
			Id:        "",
			DenomId:   params.ClassId,
			Name:      params.Name,
			URI:       params.Uri,
			UriHash:   params.UriHash,
			Data:      params.Data,
			Sender:    class.Owner,
			Recipient: params.Recipient,
		}
		baseTx := n.base.CreateBaseTx(class.Owner, "")
		signedData, _, err := n.base.BuildAndSign(sdktype.Msgs{&createNft}, baseTx)
		if err != nil {
			log.Debug("create ddc", "buildandsign error:", err.Error())
			return types.ErrBuildAndSign
		}
		msgsBytes, _ := json.Marshal(sdktype.Msgs{&createNft})

		ttx := models.TTX{
			ProjectID:     params.ProjectID,
			Hash:          res.TxHash,
			Timestamp:     null.Time{Time: time.Now()},
			Message:       null.JSONFrom(msgsBytes),
			Sender:        null.StringFrom(class.Owner),
			TaskID:        null.StringFrom(taskId),
			GasUsed:       null.Int64From(0),
			OriginData:    null.BytesFromPtr(&signedData),
			OperationType: models.TTXSOperationTypeMintNFT,
			Status:        models.TTXSStatusUndo,
			Tag:           null.JSONFrom(params.Tag),
			Retry:         null.Int8From(0),
		}
		err = ttx.Insert(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Error("create ddc", "ttx insert error: ", err)
			return types.ErrInternal
		}

		//class locked
		class.Status = models.TDDCClassesStatusPending
		class.LockedBy = null.Uint64FromPtr(&ttx.ID)
		ok, err := class.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Error("create ddc", "class status update error: ", err)
			return types.ErrInternal
		}
		if ok != 1 {
			log.Error("create ddc", "class status update error: ", err)
			return types.ErrInternal
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &dto.TxRes{TaskId: taskId}, nil
}

func (n Nft) Show(params dto.NftByNftIdP) (*dto.NftR, error) {
	panic("implement me")
}

func (n Nft) Update(params dto.EditNftByNftIdP) (*dto.TxRes, error) {
	// ValidateSigner
	if err := n.base.ValidateSigner(params.Sender, params.ProjectID); err != nil {
		return nil, err
	}
	tDDC, err := models.TDDCNFTS(models.TDDCNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TDDCNFTWhere.ClassID.EQ(params.ClassId),
		models.TDDCNFTWhere.NFTID.EQ(params.NftId),
		models.TDDCNFTWhere.Owner.EQ(params.Sender)).
		One(context.Background(), boil.GetContextDB())

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), service.SqlNotFound) {
			//404
			return nil, types.ErrNotFound
		}
		//500
		log.Error("edit ddc by nftId", "query ddc error:", err.Error())
		return nil, types.ErrInternal
	}
	//404
	if tDDC.Status == models.TDDCNFTSStatusBurned {
		return nil, types.ErrNotFound
	}
	//400
	if tDDC.Status != models.TDDCNFTSStatusActive {
		return nil, types.ErrNftStatus
	}

	//非必填保留数据
	uri := params.Uri
	if uri == "" {
		uri = tDDC.URI.String
	}
	data := params.Data
	if data == "" {
		data = tDDC.Metadata.String
	}

	// create rawMsg
	msgEditNFT := nft.MsgEditNFT{
		Id:      tDDC.NFTID,
		DenomId: tDDC.ClassID,
		Name:    params.Name,
		URI:     uri,
		Data:    data,
		Sender:  params.Sender,
		UriHash: "[do-not-modify]",
	}

	// build and sign transaction
	baseTx := n.base.CreateBaseTx(params.Sender, "")
	signedData, txHash, err := n.base.BuildAndSign(sdktype.Msgs{&msgEditNFT}, baseTx)
	if err != nil {
		log.Debug("edit nft by nftId", "BuildAndSign error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		// Tx into database
		messageByte, _ := json.Marshal(msgEditNFT)
		code := fmt.Sprintf("%s%s%s", params.Sender, models.TTXSOperationTypeEditNFT, time.Now().String())
		taskId = n.base.EncodeData(code)

		// Tx into database
		txId, err := n.base.UndoTxIntoDataBase(
			params.Sender,
			models.TTXSOperationTypeEditNFT,
			taskId, txHash,
			params.ProjectID,
			signedData,
			messageByte,
			params.Tag,
			int64(baseTx.Gas), exec)
		if err != nil {
			log.Debug("edit ddc by nftId", "Tx into database error:", err.Error())
			return err
		}

		// lock the NFT
		tDDC.Status = models.TDDCNFTSStatusPending
		tDDC.LockedBy = null.Uint64From(txId)
		ok, err := tDDC.Update(context.Background(), exec, boil.Infer())
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

func (n Nft) Delete(params dto.DeleteNftByNftIdP) (*dto.TxRes, error) {
	panic("implement me")
}

func (n Nft) History(params dto.NftOperationHistoryByNftIdP) (*dto.BNftOperationHistoryByNftIdRes, error) {
	result := &dto.BNftOperationHistoryByNftIdRes{
		PageRes: dto.PageRes{
			Offset:     params.Offset,
			Limit:      params.Limit,
			TotalCount: 0,
		},
		OperationRecords: []*dto.OperationRecord{},
	}
	res, err := models.TDDCNFTS(
		models.TDDCNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TDDCNFTWhere.ClassID.EQ(params.ClassID),
		models.TDDCNFTWhere.NFTID.EQ(params.NftId),
	).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), service.SqlNotFound)) {
		//404
		return nil, types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("query ddc operation history", "query ddc error:", err.Error())
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
		if strings.Contains(err.Error(), service.SqlNotFound) {
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
