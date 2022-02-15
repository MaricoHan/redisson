package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"strings"

	sdk "github.com/irisnet/core-sdk-go"
	"github.com/irisnet/core-sdk-go/bank"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type Base struct {
	sdkClient sdk.Client
	gas       uint64
	coins     sdktype.DecCoins
}

func SqlNoFound() string {
	return "records not exist"
}

func NewBase(sdkClient sdk.Client, gas uint64, denom string, amount int64) *Base {
	return &Base{
		sdkClient: sdkClient,
		gas:       gas,
		coins:     sdktype.NewDecCoins(sdktype.NewDecCoin(denom, sdktype.NewInt(amount))),
	}
}

func (m Base) CreateBaseTx(keyName, keyPassword string) sdktype.BaseTx {
	//from := "t_" + keyName
	return sdktype.BaseTx{
		From:     keyName,
		Gas:      m.gas,
		Fee:      m.coins,
		Mode:     sdktype.Commit,
		Password: keyPassword,
	}
}

func (m Base) BuildAndSign(msgs sdktype.Msgs, baseTx sdktype.BaseTx) ([]byte, string, error) {
	sigData, err := m.sdkClient.BuildAndSign(msgs, baseTx)
	if err != nil {
		return nil, "", err
	}
	hashBz := sha256.Sum256(sigData)
	hash := strings.ToUpper(hex.EncodeToString(hashBz[:]))
	return sigData, hash, nil
}

func (m Base) BuildAndSend(msgs sdktype.Msgs, baseTx sdktype.BaseTx) (sdktype.ResultTx, error) {
	sigData, err := m.sdkClient.BuildAndSend(msgs, baseTx)
	if err != nil {
		return sigData, err
	}
	return sigData, nil
}

// TxIntoDataBase operationType : issue_class,mint_nft,edit_nft,edit_nft_batch,burn_nft,burn_nft_batch
func (m Base) TxIntoDataBase(AppID uint64, txHash string, signedData []byte, operationType string, status string, exec boil.ContextExecutor) (uint64, error) {
	// Tx into database
	ttx := models.TTX{
		AppID:         AppID,
		Hash:          txHash,
		OriginData:    null.BytesFrom(signedData),
		OperationType: operationType,
		Status:        status,
	}
	err := ttx.Insert(context.Background(), exec, boil.Infer())
	if err != nil {
		return 0, err
	}
	return ttx.ID, err
}

// ValidateTx validate tx status
func (m Base) ValidateTx(hash string) (*models.TTX, error) {
	tx, err := models.TTXS(models.TTXWhere.Hash.EQ(hash)).OneG(context.Background())
	if err == sql.ErrNoRows || strings.Contains(err.Error(), SqlNoFound()) {
		return tx, nil
	} else if err != nil {
		return tx, err
	}

	// pending
	if tx.Status == models.TTXSStatusPending {
		return tx, types.ErrTXStatusPending
	}

	// undo
	if tx.Status == models.TTXSStatusUndo {
		return tx, types.ErrTXStatusUndo
	}

	// success
	if tx.Status == models.TTXSStatusSuccess {
		return tx, types.ErrTXStatusSuccess
	}

	return tx, nil
}

func (m Base) CreateGasMsg(inputAddress string, outputAddress []string) bank.MsgMultiSend {
	accountGas := config.Get().Chain.AccoutGas
	denom := config.Get().Chain.Denom
	inputCoins := sdktype.NewCoins(sdktype.NewCoin(denom, sdktype.NewInt(accountGas*int64(len(outputAddress)))))
	outputCoins := sdktype.NewCoins(sdktype.NewCoin(denom, sdktype.NewInt(accountGas)))
	inputs := []bank.Input{
		{
			Address: inputAddress,
			Coins:   inputCoins,
		},
	}
	var outputs []bank.Output
	for _, v := range outputAddress {
		outputs = append(outputs, bank.Output{
			Address: v,
			Coins:   outputCoins,
		})
	}
	msg := bank.MsgMultiSend{
		Inputs:  inputs,
		Outputs: outputs,
	}
	return msg
}
func (m Base) calculateGas(origindata []byte) uint64 {
	return config.Get().Chain.GasCoefficient * uint64(len(origindata))
}
