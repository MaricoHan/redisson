package wenchangchain_ddc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

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
)

const nftp = "nftp"

type Nft struct {
	base *service.Base
}

func NewDDCNFT(base *service.Base) *service.NFTBase {
	return &service.NFTBase{
		Module: service.DDC,
		Service: &Nft{
			base: base,
		},
	}
}

func (n Nft) List(params dto.NftsP) (*dto.NftsRes, error) {
	panic("implement me")
}

func (n Nft) Create(params dto.CreateNftsP) (*dto.TxRes, error) {
	DDC721Service := service.DDCClient.GetDDC721Service(true)
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
		if err := n.base.ValidateDDCSigner(class.Owner, params.ProjectID); err != nil {
			return err
		}

		if params.Uri == "" {
			params.Uri = "-"
		}

		if params.Recipient == "" {
			//默认为 NFT 类别的权属者地址
			params.Recipient = class.Owner
		}

		// ValidateRecipient
		if err := n.base.ValidateDDCRecipient(params.Recipient, params.ProjectID); err != nil {
			return err
		}

		//taskId
		code := fmt.Sprintf("%s%s%s", class.Owner, models.TDDCTXSOperationTypeMintNFT, time.Now().String())
		taskId = n.base.EncodeData(code)

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

		//签名后的交易计算动态 gas
		addr, err := DDC721Service.Bech32ToHex(params.Recipient)
		if err != nil {
			//ddc sdk 自定义 err
			return err
		}
		res, err := DDC721Service.SafeMint(&bind.TransactOpts{}, addr, params.Uri, []byte(params.Data))
		if err != nil {
			log.Error("create ddc", "get hash and gasLimit error:", err.Error())
			return types.ErrInternal
		}
		err = n.base.GasThan(params.ChainID, res.GasLimit, params.PlatFormID)
		if err != nil {
			return types.NewAppError(types.RootCodeSpace, types.ErrGasNotEnough, err.Error())
		}

		msgsBytes, _ := json.Marshal(sdktype.Msgs{&createNft})

		//null.Int64From(int64(res.GasLimit))
		ttx := models.TDDCTX{
			ProjectID:     params.ProjectID,
			Hash:          taskId,
			Timestamp:     null.Time{Time: time.Now()},
			Message:       null.JSONFrom(msgsBytes),
			Sender:        null.StringFrom(class.Owner),
			TaskID:        null.StringFrom(taskId),
			GasUsed:       null.Int64From(int64(res.GasLimit)),
			OriginData:    null.BytesFromPtr(&msgsBytes),
			OperationType: models.TDDCTXSOperationTypeMintNFT,
			Status:        models.TDDCTXSStatusUndo,
			Tag:           null.JSONFrom(params.Tag),
			Retry:         null.Int8From(0),
			BizFee:        null.Int64From(service.MintFee),
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
	DDC721Service := service.DDCClient.GetDDC721Service(true)
	var taskId string
	err := modext.Transaction(func(exec boil.ContextExecutor) error {
		// ValidateSigner
		if err := n.base.ValidateDDCSigner(params.Sender, params.ProjectID); err != nil {
			return err
		}

		//query ddc
		tDDC, err := models.TDDCNFTS(models.TDDCNFTWhere.ProjectID.EQ(params.ProjectID),
			models.TDDCNFTWhere.ClassID.EQ(params.ClassId),
			models.TDDCNFTWhere.NFTID.EQ(params.NftId),
			models.TDDCNFTWhere.Owner.EQ(params.Sender)).
			One(context.Background(), boil.GetContextDB())
		if err != nil {
			if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), service.SqlNotFound) {
				//404
				return types.ErrNotFound
			}
			//500
			log.Error("edit ddc by nftId", "query ddc error:", err.Error())
			return types.ErrInternal
		}

		//404
		if tDDC.Status == models.TDDCNFTSStatusBurned {
			return types.ErrNotFound
		}
		//400
		if tDDC.Status != models.TDDCNFTSStatusActive {
			return types.ErrNftStatus
		}

		//非必填保留数据
		//uri 为空不上链
		uri := params.Uri
		if uri == "" {
			uri = tDDC.URI.String
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

			// Tx into database
			messageByte, _ := json.Marshal(msgEditNFT)
			code := fmt.Sprintf("%s%s%s", params.Sender, models.TDDCTXSOperationTypeEditNFT, time.Now().String())
			taskId = n.base.EncodeData(code)

			//tx 表
			ttx := models.TDDCTX{
				ProjectID:     params.ProjectID,
				Hash:          taskId,
				OriginData:    null.BytesFrom(messageByte),
				OperationType: models.TDDCTXSOperationTypeEditNFT,
				Status:        models.TDDCTXSStatusSuccess,
				Sender:        null.StringFrom(params.Sender),
				Message:       null.JSONFrom(messageByte),
				TaskID:        null.StringFrom(taskId),
				GasUsed:       null.Int64From(0),
				Tag:           null.JSONFrom(params.Tag),
				Retry:         null.Int8From(0),
			}
			err = ttx.Insert(context.Background(), exec, boil.Infer())
			if err != nil {
				log.Error("edit ddc by nftId", "tx into database error:", err.Error())
				return err
			}

			//msg 表
			tmsg := models.TMSG{
				ProjectID: params.ProjectID,
				TXHash:    taskId,
				Module:    models.TDDCMSGSModuleNFT,
				Operation: models.TDDCTXSOperationTypeEditNFT,
				Signer:    params.Sender,
				Timestamp: null.TimeFrom(time.Now()),
				Message:   messageByte,
			}
			err = tmsg.Insert(context.Background(), exec, boil.Infer())
			if err != nil {
				log.Error("edit ddc by nftId", "tx into database error:", err.Error())
				return err
			}

			//update ddc
			tDDC.Name = null.StringFrom(params.Name)
			tDDC.URI = null.StringFrom(uri)
			tDDC.Metadata = null.StringFrom(data)
			tDDC.URIHash = null.StringFrom("[do-not-modify]")
			ok, err := tDDC.Update(context.Background(), exec, boil.Infer())
			if err != nil {
				log.Error("edit ddc", "ddc status update error: ", err)
				return types.ErrInternal
			}
			if ok != 1 {
				log.Error("edit ddc", "ddc status update error: ", err)
				return types.ErrInternal
			}
		} else { //uri 不为空需要上链
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

			//签名后的交易计算动态 gas
			ddcId, err := strconv.ParseInt(params.NftId, 10, 64)
			if err != nil {
				log.Error("edit ddc by nftId", "convert ddcId error:", err.Error())
				return types.ErrInternal
			}
			res, err := DDC721Service.SetURI(&bind.TransactOpts{}, ddcId, params.Uri)
			if err != nil {
				log.Error("edit ddc by nftId", "get hash and gasLimit error:", err.Error())
				return types.ErrInternal
			}
			err = n.base.GasThan(params.ChainID, res.GasLimit, params.PlatFormID)
			if err != nil {
				return types.NewAppError(types.RootCodeSpace, types.ErrGasNotEnough, err.Error())
			}

			// Tx into database
			messageByte, _ := json.Marshal(msgEditNFT)
			code := fmt.Sprintf("%s%s%s", params.Sender, models.TTXSOperationTypeEditNFT, time.Now().String())
			taskId = n.base.EncodeData(code)

			// Tx into database
			txId, err := n.base.UndoDDCTxIntoDataBase(
				params.Sender,
				models.TTXSOperationTypeEditNFT,
				taskId, res.TxHash,
				params.ProjectID,
				messageByte,
				messageByte,
				params.Tag,
				int64(res.GasLimit), exec)
			if err != nil {
				log.Error("edit ddc by nftId", "tx into database error:", err.Error())
				return err
			}

			// locked by txId
			tDDC.Status = models.TDDCNFTSStatusPending
			tDDC.LockedBy = null.Uint64From(txId)
			ok, err := tDDC.Update(context.Background(), exec, boil.Infer())
			if err != nil {
				log.Error("edit ddc", "ddc status update error: ", err)
				return types.ErrInternal
			}
			if ok != 1 {
				log.Error("edit ddc", "ddc status update error: ", err)
				return types.ErrInternal
			}
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
		qm.From(models.TableNames.TDDCMSGS),
		qm.Select(models.TDDCMSGColumns.TXHash,
			models.TDDCMSGColumns.Operation,
			models.TDDCMSGColumns.Signer,
			models.TDDCMSGColumns.Recipient,
			models.TDDCMSGColumns.Timestamp),
		models.TDDCMSGWhere.ProjectID.EQ(params.ProjectID),
		models.TDDCMSGWhere.NFTID.EQ(null.StringFrom(res.NFTID)),
	}
	if params.Txhash != "" {
		queryMod = append(queryMod, models.TDDCMSGWhere.TXHash.EQ(params.Txhash))
	}
	if params.Signer != "" {
		queryMod = append(queryMod, models.TDDCMSGWhere.Signer.EQ(params.Signer))
	}
	if params.Operation != "" {
		queryMod = append(queryMod, models.TDDCMSGWhere.Operation.EQ(params.Operation))
	}
	if params.StartDate != nil {
		queryMod = append(queryMod, models.TDDCMSGWhere.CreateAt.GTE(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TDDCMSGWhere.CreateAt.LTE(*params.EndDate))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TDDCMSGWhere.CreateAt)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TDDCMSGWhere.CreateAt)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}
	var modelResults []*models.TDDCMSG
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
