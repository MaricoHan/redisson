package wenchangchain_ddc

import (
	"context"
	"database/sql"
	"fmt"
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
	"strings"
	"time"
)

type DDCNftTransfer struct {
	base *service.Base
}

func NewDDCNftTransfer(base *service.Base) *service.TransferBase {
	return &service.TransferBase{
		Module:  service.DDC,
		Service: &DDCNftTransfer{base: base},
	}
}

func (D DDCNftTransfer) TransferNFTClass(params dto.TransferNftClassByIDP) (*dto.TxRes, error) {
	//不能自己转让给自己
	//400
	if params.Recipient == params.Owner {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSelfTransfer)
	}

	// ValidateSigner
	if err := D.base.ValidateSigner(params.Owner, params.ProjectID); err != nil {
		return nil, err
	}

	// ValidateRecipient
	if err := D.base.ValidateRecipient(params.Recipient, params.ProjectID); err != nil {
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
		taskId = D.base.EncodeData(code)
		hash := D.base.EncodeData(string(msgsByte))
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

func (D DDCNftTransfer) TransferNFT(params dto.TransferNftByNftIdP) (*dto.TxRes, error) {
	return nil, nil
}
