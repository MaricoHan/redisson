package wenchangchain_ddc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	service2 "github.com/bianjieai/ddc-sdk-go/app/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

type DDC struct {
	base          *service.Base
	ddc721Service *service2.DDC721Service
}

func NewDDC(base *service.Base) *service.NFTBase {
	client := service.NewDDCClient()
	ddc721Service := client.GetDDC721Service(true)

	return &service.NFTBase{
		Module: service.DDC,
		Service: &DDC{
			base:          base,
			ddc721Service: ddc721Service,
		},
	}
}

func (d DDC) Create(params dto.CreateNftsP) (*dto.TxRes, error) {
	// 若接收者地址不为空，则校验其格式；否则在service中将其默认设为NFT类别的权属者地址
	if params.Recipient != "" {
		// 校验接收者地址是否满足当前链的地址规范
		if !common.IsHexAddress(params.Recipient) {
			return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
		}
	}
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
		if err := d.base.ValidateDDCSigner(class.Owner, params.ProjectID); err != nil {
			return err
		}

		//if params.Uri == "" {
		//	params.Uri = "-"
		//}
		if params.Recipient == "" {
			//默认为 NFT 类别的权属者地址
			params.Recipient = class.Owner
		}

		// ValidateRecipient
		if err := d.base.ValidateDDCRecipient(params.Recipient, params.ProjectID); err != nil {
			return err
		}

		//taskId
		code := fmt.Sprintf("%s%s%s", class.Owner, models.TDDCTXSOperationTypeMintNFT, time.Now().String())
		taskId = d.base.EncodeData(code)

		//platform 发行
		platform, err := models.TDDCAccounts(
			models.TDDCAccountWhere.ProjectID.EQ(uint64(rootProjectID)),
			models.TDDCAccountWhere.ID.EQ(uint64(platformIDInTable)),
		).OneG(context.Background())
		if err != nil {
			//500
			log.Error("create ddc", "query platform error:", err.Error())
			return types.ErrInternal
		}

		//离线数据组装为 msg 入表传到 block sync
		createNft := nft.MsgMintNFT{
			Id:        "",
			DenomId:   params.ClassId,
			Name:      params.Name,
			URI:       params.Uri,
			UriHash:   params.UriHash,
			Data:      params.Data,
			Sender:    platform.Address,
			Recipient: params.Recipient,
		}

		//签名后的交易计算动态 gas
		opts := &bind.TransactOpts{
			From: common.HexToAddress(platform.Address),
		}

		res, err := d.ddc721Service.SafeMint(opts, params.Recipient, params.Uri, []byte(params.Data))
		if err != nil {
			log.Error("create ddc", "get hash and gasLimit error:", err.Error())
			return types.ErrInternal
		}
		err = d.base.GasThan(params.ChainID, res.GasLimit, params.PlatFormID)
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
		if err != nil || ok != 1 {
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

func (d DDC) Show(params dto.NftByNftIdP) (*dto.NftR, error) {
	//查出ddc
	tDDC, err := models.TDDCNFTS(models.TDDCNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TDDCNFTWhere.ClassID.EQ(params.ClassId),
		models.TDDCNFTWhere.NFTID.EQ(params.NftId)).
		One(context.Background(), boil.GetContextDB())
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), service.SqlNotFound) {
			//404
			return nil, types.ErrNotFound
		}
		//500
		log.Error("ddc by ddcId", "query ddc error:", err.Error())
		return nil, types.ErrInternal
	}
	//检验ddc状态
	if tDDC.Status == models.TDDCNFTSStatusPending {
		return nil, types.ErrNftStatus
	}
	//查出class
	class, err := models.TDDCClasses(models.TDDCClassWhere.ClassID.EQ(params.ClassId)).One(context.Background(), boil.GetContextDB())
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), service.SqlNotFound) {
			//404
			return nil, types.ErrNotFound
		}
		//500
		log.Error("ddc by ddcId", "query ddc class error:", err.Error())
		return nil, types.ErrInternal
	}

	result := &dto.NftR{
		Id:          tDDC.NFTID,
		Name:        tDDC.Name.String,
		ClassId:     tDDC.ClassID,
		ClassName:   class.Name.String,
		ClassSymbol: class.Symbol.String,
		Uri:         tDDC.URI.String,
		UriHash:     tDDC.URIHash.String,
		Data:        tDDC.Metadata.String,
		Owner:       tDDC.Owner,
		Status:      tDDC.Status,
		TxHash:      tDDC.TXHash,
		Timestamp:   tDDC.Timestamp.Time.String(),
	}

	return result, nil
}

func (d DDC) Update(params dto.EditNftByNftIdP) (*dto.TxRes, error) {
	var taskId string
	err := modext.Transaction(func(exec boil.ContextExecutor) error {
		// ValidateSigner
		if err := d.base.ValidateDDCSigner(params.Sender, params.ProjectID); err != nil {
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

			uriHash := tDDC.URIHash.String

			// create rawMsg
			msgEditNFT := nft.MsgEditNFT{
				Id:      tDDC.NFTID,
				DenomId: tDDC.ClassID,
				Name:    params.Name,
				URI:     uri,
				Data:    data,
				Sender:  params.Sender,
				UriHash: uriHash,
			}

			// Tx into database
			messageByte, _ := json.Marshal(msgEditNFT)
			code := fmt.Sprintf("%s%s%s", params.Sender, models.TDDCTXSOperationTypeEditNFT, time.Now().String())
			taskId = d.base.EncodeData(code)
			hash := "0x" + d.base.EncodeData(string(messageByte)+time.Now().String())
			//tx 表
			ttx := models.TDDCTX{
				ProjectID:     params.ProjectID,
				Hash:          hash,
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
			tmsg := models.TDDCMSG{
				ProjectID: params.ProjectID,
				TXHash:    hash,
				Module:    models.TDDCMSGSModuleNFT,
				Operation: models.TDDCMSGSOperationEdit,
				Signer:    params.Sender,
				Timestamp: null.TimeFrom(time.Now()),
				Message:   messageByte,
			}
			err = tmsg.Insert(context.Background(), exec, boil.Infer())
			if err != nil {
				log.Error("edit ddc by nftId", "msg into database error:", err.Error())
				return err
			}

			//update ddc
			tDDC.Name = null.StringFrom(params.Name)
			tDDC.URI = null.StringFrom(uri)
			tDDC.Metadata = null.StringFrom(data)
			ok, err := tDDC.Update(context.Background(), exec, boil.Infer())
			if ok != 1 {
				return types.ErrNftStatus
			}
			if err != nil {
				log.Error("edit ddc", "ddc status update error: ", err)
				return types.ErrInternal
			}
		} else { //uri 不为空判断链上 uri 是否为空
			ddcID, err := strconv.ParseInt(params.NftId, 10, 64)
			if err != nil {
				log.Error("edit ddc", "convert ddcID error: ", err)
				return types.ErrInternal
			}
			uriChain, err := d.ddc721Service.DdcURI(ddcID)
			if err != nil {
				log.Error("edit ddc", "query ddc uri error: ", err)
				return types.ErrInternal
			}
			//链上 uri 为空,上链
			if uriChain == "" {
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
					UriHash: tDDC.URIHash.String,
				}

				//签名后的交易计算动态 gas
				ddcId, err := strconv.ParseInt(params.NftId, 10, 64)
				if err != nil {
					log.Error("edit ddc by nftId", "convert ddcId error:", err.Error())
					return types.ErrInternal
				}
				opts := &bind.TransactOpts{
					From: common.HexToAddress(params.Sender),
				}

				res, err := d.ddc721Service.SetURI(opts, ddcId, params.Uri)
				if err != nil {
					log.Error("edit ddc by nftId", "get hash and gasLimit error:", err.Error())
					return types.ErrInternal
				}
				err = d.base.GasThan(params.ChainID, res.GasLimit, params.PlatFormID)
				if err != nil {
					return types.NewAppError(types.RootCodeSpace, types.ErrGasNotEnough, err.Error())
				}

				// Tx into database
				messageByte, _ := json.Marshal(msgEditNFT)
				code := fmt.Sprintf("%s%s%s", params.Sender, models.TTXSOperationTypeEditNFT, time.Now().String())
				taskId = d.base.EncodeData(code)

				// Tx into database
				txId, err := d.base.UndoDDCTxIntoDataBase(
					params.Sender,
					models.TTXSOperationTypeEditNFT,
					taskId,
					res.TxHash,
					params.ProjectID,
					messageByte,
					params.Tag,
					int64(res.GasLimit),
					0,
					exec)
				if err != nil {
					log.Error("edit ddc by nftId", "tx into database error:", err.Error())
					return err
				}

				// locked by txId
				tDDC.Status = models.TDDCNFTSStatusPending
				tDDC.LockedBy = null.Uint64From(txId)
				ok, err := tDDC.Update(context.Background(), exec, boil.Infer())
				if err != nil || ok != 1 {
					log.Error("edit ddc", "ddc status update error: ", err)
					return types.ErrInternal
				}
			} else {
				return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrUriChain)
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &dto.TxRes{TaskId: taskId}, nil
}

func (d DDC) Delete(params dto.DeleteNftByNftIdP) (*dto.TxRes, error) {
	//ValidateSigner
	if err := d.base.ValidateDDCSigner(params.Sender, params.ProjectID); err != nil {
		return nil, err
	}
	//查出要删除的ddc
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
		log.Error("delete ddc by ddcId", "query nft error:", err.Error())
		return nil, types.ErrInternal
	}
	//判断ddc的状态
	// 404
	if tDDC.Status == models.TDDCNFTSStatusBurned {
		return nil, types.ErrNotFound
	}
	//400
	if tDDC.Status != models.TDDCNFTSStatusActive {
		return nil, types.ErrNftStatus
	}

	//组装rawMsg
	msgBurnNFT := nft.MsgBurnNFT{
		Id:      tDDC.NFTID,
		DenomId: tDDC.ClassID,
		Sender:  params.Sender,
	}
	messageByte, _ := json.Marshal(msgBurnNFT)
	//生成taskId
	code := fmt.Sprintf("%s%s%s", params.Sender, models.TDDCTXSOperationTypeBurnNFT, time.Now().String())
	taskId := d.base.EncodeData(code)

	//获取gasLimit和txHash
	opts := bind.TransactOpts{
		From: common.HexToAddress(params.Sender),
	}
	ddcId, _ := strconv.ParseInt(tDDC.NFTID, 10, 64)
	res, err := d.ddc721Service.Burn(&opts, ddcId)
	if err != nil {
		log.Error("delete ddc by ddcId", "failed to get gasLimit and txHash", err.Error())
		return nil, types.ErrInternal
	}

	//tx存数据库
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		// Tx into database
		txId, err := d.base.UndoDDCTxIntoDataBase(params.Sender,
			models.TDDCTXSOperationTypeBurnNFT,
			taskId,
			taskId,
			params.ProjectID,
			messageByte,
			params.Tag,
			int64(res.GasLimit),
			service.BurnFee,
			exec)
		if err != nil {
			log.Error("delete ddc by ddcId", "tx into database error:", err.Error())
			return types.ErrInternal
		}

		// lock the NFT
		tDDC.Status = models.TDDCNFTSStatusPending
		tDDC.LockedBy = null.Uint64From(txId)
		ok, err := tDDC.Update(context.Background(), exec, boil.Infer())
		if err != nil || ok != 1 {
			log.Error("delete ddc by ddcId", "failed to lock the nft:", err.Error())
			return types.ErrInternal
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return &dto.TxRes{TaskId: taskId}, nil
}
func (d DDC) List(params dto.NftsP) (*dto.NftsRes, error) {
	result := &dto.NftsRes{}
	var err error
	result.Offset = params.Offset
	result.Limit = params.Limit
	result.Nfts = []*dto.Nft{}
	//组装查询条件queryMod
	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TDDCNFTS),
		models.TDDCNFTWhere.ProjectID.EQ(params.ProjectID),
	}
	if params.Id != "" {
		queryMod = append(queryMod, models.TDDCNFTWhere.NFTID.EQ(params.Id))
	}
	if params.ClassId != "" {
		queryMod = append(queryMod, models.TDDCNFTWhere.ClassID.EQ(params.ClassId))
	}
	if params.Owner != "" {
		queryMod = append(queryMod, models.TDDCNFTWhere.Owner.EQ(params.Owner))
	}
	if params.TxHash != "" {
		queryMod = append(queryMod, models.TDDCNFTWhere.TXHash.EQ(params.TxHash))
	}
	if params.Status != "" {
		queryMod = append(queryMod, models.TDDCNFTWhere.Status.EQ(params.Status))
	}
	if params.StartDate != nil {
		queryMod = append(queryMod, models.TDDCNFTWhere.Timestamp.GTE(null.TimeFromPtr(params.StartDate)))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TDDCNFTWhere.Timestamp.LTE(null.TimeFromPtr(params.EndDate)))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "ID_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TDDCNFTColumns.NFTID)
		case "ID_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TDDCNFTColumns.NFTID)
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TDDCNFTColumns.Timestamp)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TDDCNFTColumns.Timestamp)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}
	//查询nft和对应的class信息
	var modelResults []*models.TDDCNFT
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

		//统计查询结果中有多少种不同的classId，并存入classIds
		tempMap := map[string]byte{}
		for _, m := range modelResults {
			tempMap[m.ClassID] = 0
		}
		for k, _ := range tempMap {
			classIds = append(classIds, k)
		}
		//查询所有涉及到的class的相关信息，并存入classByIds
		qMod := []qm.QueryMod{
			qm.From(models.TableNames.TDDCClasses),
			qm.Select(models.TDDCClassColumns.ClassID, models.TDDCClassColumns.Name, models.TDDCClassColumns.Symbol),
			models.TDDCClassWhere.ClassID.IN(classIds),
		}
		err = models.NewQuery(qMod...).Bind(context.Background(), exec, &classByIds)
		return err
	})

	if err != nil {
		// records not exist
		if strings.Contains(err.Error(), service.SqlNotFound) {
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
		//找到该nft对应的class的信息
		for _, class := range classByIds {
			if class.ClassId == modelResult.ClassID {
				nft.ClassName = class.Name
				nft.ClassSymbol = class.Symbol
				break
			}
		}
		nfts = append(nfts, nft)
	}
	result.Nfts = nfts
	return result, nil

}
func (d DDC) History(params dto.NftOperationHistoryByNftIdP) (*dto.BNftOperationHistoryByNftIdRes, error) {
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
		queryMod = append(queryMod, models.TDDCMSGWhere.Timestamp.GTE(null.TimeFromPtr(params.StartDate)))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TDDCMSGWhere.Timestamp.LTE(null.TimeFromPtr(params.EndDate)))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TDDCMSGColumns.Timestamp)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TDDCMSGColumns.Timestamp)
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
