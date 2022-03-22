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
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
)

type DDC721Transfer struct {
	base          *service.Base
	ddc721Service *service2.DDC721Service
}

func NewDDCTransfer(base *service.Base) *service.TransferBase {
	client := service.NewDDCClient()
	ddc721Service := client.GetDDC721Service(true)
	return &service.TransferBase{
		Module: service.DDC,
		Service: &DDC721Transfer{
			base:          base,
			ddc721Service: ddc721Service,
		},
	}
}
func (d DDC721Transfer) TransferNFTClass(params dto.TransferNftClassByIDP) (*dto.TxRes, error) {
	// 校验接收者地址是否满足当前链的地址规范
	if !common.IsHexAddress(params.Recipient) {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
	}

	//不能自己转让给自己
	//400
	if params.Recipient == params.Owner {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSelfTransfer)
	}

	// ValidateSigner
	if err := d.base.ValidateDDCSigner(params.Owner, params.ProjectID); err != nil {
		return nil, err
	}

	// ValidateRecipient
	if err := d.base.ValidateDDCRecipient(params.Recipient, params.ProjectID); err != nil {
		return nil, err
	}
	//判断class
	class, err := models.TDDCClasses(
		models.TDDCClassWhere.ClassID.EQ(params.ClassID),
		models.TDDCClassWhere.ProjectID.EQ(params.ProjectID),
		models.TDDCClassWhere.Owner.EQ(params.Owner)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), service.SqlNotFound)) {
		//404
		return nil, types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("transfer ddc class", "query class error:", err.Error())
		return nil, types.ErrInternal
	}

	if class.Status != models.TDDCClassesStatusActive {
		//400
		return nil, types.ErrNftClassStatus
	}

	//msg
	msgs := nft.MsgTransferDenom{
		Id:        params.ClassID,
		Sender:    params.Owner,
		Recipient: params.Recipient,
	}

	//sign
	msgsByte, err := msgs.Marshal()
	if err != nil {
		log.Debug("transfer ddc class", "msgs Marshal error:", err.Error())
		return nil, types.ErrBuildAndSign
	}

	var taskId string
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		code := fmt.Sprintf("%s%s%s", params.Owner, models.TDDCTXSOperationTypeTransferClass, time.Now().String())
		taskId = d.base.EncodeData(code)
		hash := d.base.EncodeData(string(msgsByte))
		// Tx into database
		ttx := models.TDDCTX{
			ProjectID:     params.ProjectID,
			Hash:          hash,
			OriginData:    null.BytesFrom(msgsByte),
			OperationType: models.TDDCTXSOperationTypeTransferClass,
			Status:        models.TDDCTXSStatusSuccess,
			Sender:        null.StringFrom(params.Owner),
			Message:       null.JSONFrom(msgsByte),
			TaskID:        null.StringFrom(taskId),
			GasUsed:       null.Int64From(0),
			Tag:           null.JSONFrom(params.Tag),
			Retry:         null.Int8From(0),
			BizFee:        null.Int64From(0),
		}
		err = ttx.Insert(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Debug("transfer nft class", "Tx Into DataBase error:", err.Error())
			return err
		}

		msgsModule := models.TDDCMSG{
			ProjectID: ttx.ProjectID,
			TXHash:    ttx.Hash,
			Module:    models.TDDCMSGSModuleNFT,
			Operation: models.TDDCMSGSOperationTransferClass,
			Signer:    ttx.Sender.String,
			Recipient: null.StringFrom(params.Recipient),
			Timestamp: null.TimeFrom(time.Now()),
			Message:   msgsByte,
		}
		err = msgsModule.Insert(context.Background(), exec, boil.Infer())
		if err != nil {
			log.Debug("transfer nft class", "msg error:", err.Error())
			return err
		}
		class.Owner = params.Recipient
		//class status = pending && lockby = txs.id
		//class.Status = models.TDDCClassesStatusPending
		//class.LockedBy = null.Uint64From(ttx.ID)
		ok, err := class.Update(context.Background(), exec, boil.Infer())
		if err != nil {
			//500
			log.Error("transfer ddc class", "update ddc class error:", err.Error())
			return types.ErrInternal
		}
		if ok != 1 {
			log.Error("transfer ddc class", "update ddc class error:", err.Error())
			return types.ErrInternal
		}

		return err
	})
	if err != nil {
		return nil, err
	}
	return &dto.TxRes{TaskId: taskId}, nil
}
func (d DDC721Transfer) TransferNFT(params dto.TransferNftByNftIdP) (*dto.TxRes, error) {
	// 校验接收者地址是否满足当前链的地址规范
	if !common.IsHexAddress(params.Recipient) {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
	}
	// ValidateSigner
	if err := d.base.ValidateDDCSigner(params.Sender, params.ProjectID); err != nil {
		return nil, err
	}
	// ValidateRecipient
	if err := d.base.ValidateDDCRecipient(params.Recipient, params.ProjectID); err != nil {
		return nil, err
	}

	//查出要转让的ddc
	tDDC, err := models.TDDCNFTS(
		models.TDDCNFTWhere.NFTID.EQ(params.NftId),
		models.TDDCNFTWhere.ClassID.EQ(params.ClassID),
		models.TDDCNFTWhere.ProjectID.EQ(params.ProjectID),
		models.TDDCNFTWhere.Owner.EQ(params.Sender),
	).OneG(context.Background())
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows || strings.Contains(err.Error(), service.SqlNotFound) {
			//404
			return nil, types.ErrNotFound
		}
		//500
		log.Error("transfer nft", "query nft error:", err.Error())
		return nil, types.ErrInternal
	}
	//检验ddc状态
	//404
	if tDDC.Status == models.TDDCNFTSStatusBurned {
		return nil, types.ErrNotFound
	}
	//400
	if tDDC.Status != models.TDDCNFTSStatusActive {
		return nil, types.ErrNftStatus
	}

	//组装rawMsg
	msg := nft.MsgTransferNFT{
		Id:        tDDC.NFTID,
		DenomId:   params.ClassID,
		Name:      tDDC.Name.String,
		URI:       tDDC.URI.String,
		Data:      tDDC.Metadata.String,
		Sender:    params.Sender,
		Recipient: params.Recipient,
		UriHash:   tDDC.URIHash.String,
	}
	messageByte, _ := json.Marshal(msg)
	//生成taskId
	code := fmt.Sprintf("%s%s%s", params.Sender, models.TDDCTXSOperationTypeTransferNFT, time.Now().String())
	taskId := d.base.EncodeData(code)
	//获取gasLimit和txHash
	opts := bind.TransactOpts{
		From: common.HexToAddress(params.Sender),
	}
	ddcId, _ := strconv.ParseInt(tDDC.NFTID, 10, 64)
	res, err := d.ddc721Service.TransferFrom(&opts, params.Sender, params.Recipient, ddcId)
	if err != nil {
		log.Error("transfer ddc by ddcId", "failed to get gasLimit and txHash", err.Error())
		return nil, types.ErrInternal
	}

	//tx存数据库
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		// Tx into database
		txId, err := d.base.UndoDDCTxIntoDataBase(params.Sender,
			models.TDDCTXSOperationTypeTransferNFT,
			taskId,
			taskId,
			params.ProjectID,
			messageByte,
			params.Tag,
			int64(res.GasLimit),
			service.TransFer,
			exec)
		if err != nil {
			log.Error("transfer ddc by ddcId", "tx Into DataBase error:", err.Error())
			return types.ErrInternal
		}

		// lock the NFT
		tDDC.Status = models.TDDCNFTSStatusPending
		tDDC.LockedBy = null.Uint64From(txId)
		ok, err := tDDC.Update(context.Background(), exec, boil.Infer())
		if err != nil || ok != 1 {
			log.Error("transfer ddc by ddcId", "failed to lock the nft:", err.Error())
			return types.ErrInternal
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return &dto.TxRes{TaskId: taskId}, nil

}
