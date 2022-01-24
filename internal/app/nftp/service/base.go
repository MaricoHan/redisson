package service

import (
	"context"
	sdk "github.com/irisnet/core-sdk-go"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type Base struct {
	sdkClient sdk.Client
	gas       uint64
	coins     sdktype.DecCoins
}

func NewBase(sdkClient sdk.Client, gas uint64, denom string, amount int64) *Base {
	return &Base{
		sdkClient: sdkClient,
		gas:       gas,
		coins:     sdktype.NewDecCoins(sdktype.NewDecCoin(denom, sdktype.NewInt(amount))),
	}
}

func (m Base) CreateBaseTx(keyName, keyPassword string) sdktype.BaseTx {
	return sdktype.BaseTx{
		From:     keyName,
		Gas:      m.gas,
		Fee:      m.coins,
		Mode:     sdktype.Async,
		Password: keyPassword,
	}
}

func (m Base) BuildAndSign(msgs sdktype.Msgs, baseTx sdktype.BaseTx) ([]byte, string, error) {
	sigData, err := m.sdkClient.BuildAndSign(msgs, baseTx)
	if err != nil {
		return nil, "", err
	}
	txHash, err := m.sdkClient.BuildTxHash(msgs, baseTx)
	if err != nil {
		return nil, "", err
	}
	return sigData, txHash, nil
}

// TxIntoDataBase operationType : issue_class,mint_nft,edit_nft,edit_nft_batch,burn_nft,burn_nft_batch
func (m Base) TxIntoDataBase(AppID uint64, txHash string, signedData []byte, operationType string, status string) (uint64, error) {

	// Tx into database
	ttx := models.TTX{
		AppID:         AppID,
		Hash:          txHash,
		OriginData:    null.BytesFrom(signedData),
		OperationType: operationType,
		Status:        status,
	}
	err := ttx.InsertG(context.Background(), boil.Infer())

	tx, err := models.TTXS(qm.Where("hash=?", txHash)).One(context.Background(), boil.GetContextDB())

	return tx.ID, err
}
