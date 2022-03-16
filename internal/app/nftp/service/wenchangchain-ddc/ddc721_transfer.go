package wenchangchain_ddc

import (
	"context"
	"database/sql"
	"encoding/json"
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

type DDC721Transfer struct {
	Base
}

func NEWDDC721Transfer(base *service.Base) *service.TransferBase {
	return &service.TransferBase{
		Module: service.DDC,
		Service: &DDC721Transfer{
			NewBase(base),
		},
	}
}
func (svc DDC721Transfer) TransferNFTClass(params dto.TransferNftClassByIDP) (*dto.TxRes, error) {
	panic("...")
}
func (svc DDC721Transfer) TransferNFT(params dto.TransferNftByNftIdP) (*dto.TxRes, error) {
	// ValidateSigner
	if err := svc.base.ValidateSigner(params.Sender, params.ProjectID); err != nil {
		return nil, err
	}

	// ValidateRecipient
	if err := svc.base.ValidateRecipient(params.Recipient, params.ProjectID); err != nil {
		return nil, err
	}

	//查出ddc
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

	//404
	if tDDC.Status == models.TNFTSStatusBurned {
		return nil, types.ErrNotFound
	}

	//400
	if tDDC.Status != models.TNFTSStatusActive {
		return nil, types.ErrNftStatus
	}
	//组装交易
	msgs := nft.MsgTransferNFT{
		Id:        tDDC.NFTID,
		DenomId:   params.ClassID,
		Name:      tDDC.Name.String,
		URI:       tDDC.URI.String,
		Data:      tDDC.Metadata.String,
		Sender:    params.Sender,
		Recipient: params.Recipient,
		UriHash:   tDDC.URIHash.String,
	}
	messageByte, _ := json.Marshal(msgs)
	code := fmt.Sprintf("%s%s%s", params.Sender, models.TTXSOperationTypeTransferNFT, time.Now().String())
	taskId := svc.base.EncodeData(code)
	txHash := ""

	//tx存数据库
	err = modext.Transaction(func(exec boil.ContextExecutor) error {

		// Tx into database
		txId, err := svc.UndoTxIntoDataBase(params.Sender, models.TTXSOperationTypeTransferNFT, taskId, txHash,
			params.ProjectID, messageByte, params.Tag, exec)
		if err != nil {
			log.Debug("transfer nft by index", "Tx Into DataBase error:", err.Error())
			return types.ErrInternal
		}

		// lock the NFT
		tDDC.Status = models.TNFTSStatusPending
		tDDC.LockedBy = null.Uint64From(txId)
		ok, err := tDDC.Update(context.Background(), exec, boil.Infer())
		if err != nil || ok != 1 {
			return types.ErrInternal
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return &dto.TxRes{TaskId: taskId}, nil

}
