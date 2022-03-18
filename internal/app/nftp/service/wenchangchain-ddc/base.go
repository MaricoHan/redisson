package wenchangchain_ddc

import (
	"context"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type Base struct {
	base *service.Base
}

func NewBase(base *service.Base) Base {
	return Base{base}
}

// UndoTxIntoDataBase operationType : issue_class,mint_nft,edit_nft,edit_nft_batch,burn_nft,burn_nft_batch
func (b Base) UndoTxIntoDataBase(sender, operationType, taskId, txHash string, ProjectID uint64, message, tag []byte, gasUsed, bizFee int64, exec boil.ContextExecutor) (uint64, error) {

	// Tx into database
	ttx := models.TDDCTX{
		ProjectID:     ProjectID,
		Hash:          txHash,
		OperationType: operationType,
		Status:        models.TDDCTXSStatusUndo,
		Sender:        null.StringFrom(sender),
		Message:       null.JSONFrom(message),
		TaskID:        null.StringFrom(taskId),
		Tag:           null.JSONFrom(tag),
		GasUsed:       null.Int64From(gasUsed),
		Retry:         null.Int8From(0),
		BizFee:        null.Int64From(bizFee),
	}
	err := ttx.Insert(context.Background(), exec, boil.Infer())
	if err != nil {
		return 0, err
	}
	return ttx.ID, err
}
