package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	sdk "github.com/irisnet/core-sdk-go"
	"github.com/irisnet/core-sdk-go/bank"
	"github.com/irisnet/core-sdk-go/feegrant"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"strings"
)

type Base struct {
	sdkClient sdk.Client
	gas       uint64
	coins     sdktype.DecCoins
}

func SqlNoFound() string {
	return "records not exist"
}

func QueryRootAccount() (string, *types.AppError) {
	//platform address
	classOne, err := models.TAccounts(
		models.TAccountWhere.ChainID.EQ(uint64(0)),
	).OneG(context.Background())
	if err != nil {
		//500
		log.Error("create account", "query root error:", err.Error())
		return "", types.ErrInternal
	}
	pAddress := classOne.Address
	return pAddress, nil
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
	pAddress, error := QueryRootAccount()
	if error != nil{
		return nil, "", error
	}
	baseTx.FeePayer = sdktype.AccAddress(pAddress)
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
func (m Base) TxIntoDataBase(ChainID uint64, txHash string, signedData []byte, operationType string, status string, message []byte, sender, taskId string, gas int64, exec boil.ContextExecutor) (uint64, error) {
	// Tx into database
	ttx := models.TTX{
		ChainID:       ChainID,
		Hash:          txHash,
		OriginData:    null.BytesFrom(signedData),
		OperationType: operationType,
		Status:        status,
		Sender:        null.StringFrom(sender),
		Message:       null.JSONFrom(message),
		TaskID:        null.StringFrom(taskId),
		GasUsed:       null.Int64From(gas),
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

/**
Estimated gas required to issue nft
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58048328
*/
func (m Base) mintNftsGas(originData []byte, amount uint64) uint64 {
	l := uint64(len(originData))
	if l <= types.MintMinNFTDataSize {
		return uint64(float64(types.MintMinNFTGas) * config.Get().Chain.GasCoefficient)
	}
	amount -= 1
	res := (l-types.MintMinNFTIncreaseDataSize*(amount)-types.MintMinNFTDataSize)*types.MintNFTCoefficient + types.MintMinNFTGas + types.MintMinNFTIncreaseGas*(amount)
	u := float64(res) * config.Get().Chain.GasCoefficient
	return uint64(u)
}

/**
Estimated gas required to create denom
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58048352
*/
func (m Base) createDenomGas(data []byte) uint64 {
	l := uint64(len(data))
	if l <= types.CreateMinDENOMDataSize {
		return uint64(types.CreateMinDENOMGas * config.Get().Chain.GasCoefficient)
	}
	u := (l-types.CreateMinDENOMDataSize)*types.CreateDENOMCoefficient + types.CreateMinDENOMGas
	return uint64(float64(u) * config.Get().Chain.GasCoefficient)
}

/**
Estimated gas required to transfer denom
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58048356
*/
func (m Base) transferDenomGas(class *models.TClass) uint64 {
	l := len([]byte(class.ClassID)) + len([]byte(class.Status)) + len([]byte(class.Owner)) + len([]byte(class.TXHash)) + len([]byte(string(class.ChainID))) + len([]byte(string(class.ID))) + len([]byte(string(class.Offset)))
	if class.LockedBy.Valid {
		l += len([]byte(string(class.LockedBy.Uint64)))
	}
	if class.Timestamp.Valid {
		l += len([]byte(class.Timestamp.Time.String()))
	}
	if class.URIHash.Valid {
		l += len([]byte(class.URIHash.String))
	}
	if class.URI.Valid {
		l += len([]byte(class.URI.String))
	}
	if class.Name.Valid {
		l += len([]byte(class.Name.String))
	}
	if class.Symbol.Valid {
		l += len([]byte(class.Symbol.String))
	}
	if class.Description.Valid {
		l += len([]byte(class.Description.String))
	}
	if class.Metadata.Valid {
		l += len([]byte(class.Metadata.String))
	}
	if class.Extra1.Valid {
		l += len([]byte(class.Extra1.String))
	}
	if class.Extra2.Valid {
		l += len([]byte(class.Extra2.String))
	}
	if class.UpdateAt.String() != "" {
		l += len([]byte(class.UpdateAt.String()))
	}
	if class.CreateAt.String() != "" {
		l += len([]byte(class.CreateAt.String()))
	}
	if l <= types.TransferMinDENOMDataSize {
		return uint64(types.TransferMinDENOMGas * config.Get().Chain.GasCoefficient)
	}
	res := (float64(l-types.TransferMinDENOMDataSize)*types.TransferDENOMCoefficient + types.TransferMinDENOMGas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to transfer one nft
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58048358
*/
func (m Base) transferOneNftGas(data []byte) uint64 {
	l := len(data)
	if l <= types.TransferMinNFTDataSize {
		return uint64(types.TransferMinNFTGas * config.Get().Chain.GasCoefficient)
	}
	res := float64((l-types.TransferMinNFTDataSize)*types.TransferNFTCoefficient+types.TransferMinNFTGas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to transfer more nft
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58048358
*/
func (m Base) transferNftsGas(data []byte, amount uint64) uint64 {
	l := uint64(len(data))
	if l <= types.TransferMinNFTDataSize {
		return uint64(float64(types.TransferMinNFTGas) * config.Get().Chain.GasCoefficient)
	}
	res := (l-types.TransferMinNFTIncreaseDataSize*(amount-1)-types.TransferMinNFTDataSize)*types.TransferNFTCoefficient + types.TransferMinNFTGas + types.TransferMinNFTIncreaseGas*(amount-1)
	u := float64(res) * config.Get().Chain.GasCoefficient
	return uint64(u)
}

func (m Base) lenOfNft(tNft *models.TNFT) uint64 {
	len1 := len(tNft.Status + tNft.NFTID + tNft.Owner + tNft.ClassID + tNft.TXHash + tNft.Name.String + tNft.Metadata.String + tNft.URIHash.String + tNft.URI.String)
	len2 := 4 * 8 // 4 uint64
	len3 := len(tNft.CreateAt.String() + tNft.UpdateAt.String() + tNft.Timestamp.Time.String())
	return uint64(len1 + len2 + len3)
}

/**
Estimated gas required to edit nft
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049122
*/
func (m Base) editNftGas(nftLen, signLen uint64) uint64 {
	gas := types.EditNFTBaseGas + types.EditNFTLenCoefficient*nftLen + types.EditNFTSignLenCoefficient*signLen
	res := float64(gas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to edit nfts
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049126
*/
func (m Base) editBatchNftGas(nftLen, signLen uint64) uint64 {
	gas := types.EditBatchNFTBaseGas + types.EditBatchNFTLenCoefficient*nftLen + types.EditBatchNFTSignLenCoefficient*signLen
	res := float64(gas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to delete nft
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049119
*/
func (m Base) deleteNftGas(nftLen uint64) uint64 {
	gas := types.DeleteNFTBaseGas + (nftLen-types.DeleteNFTBaseLen)*types.DeleteNFTCoefficient
	res := float64(gas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to delete nfts
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049124
*/
func (m Base) deleteBatchNftGas(nftLen, n uint64) uint64 {
	basLen := types.DeleteBatchNFTBaseLen + types.DeleteBatchNFTBaseLenCoefficient*(n-1)
	baseGas := types.DeleteBatchNFTBaseGas + types.DeleteBatchNFTBaseGasCoefficient*(n-1)
	gas := (nftLen-basLen)*types.DeleteBatchCoefficient + baseGas
	res := float64(gas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to create account
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049266
*/
func (m Base) createAccount(count int64) uint64 {
	count -= 1
	res := types.CreateAccountGas + types.CreateAccountIncreaseGas*(count)
	u := float64(res) * config.Get().Chain.GasCoefficient
	return uint64(u)
}

func (m Base) Grant(address string) (string, error) {
	pAddress, error := QueryRootAccount()
	if error != nil{
		return "",error
	}
	granter, _ := sdktype.AccAddressFromBech32(pAddress)
	grantee, _ := sdktype.AccAddressFromBech32(address)

	var grant feegrant.FeeAllowanceI

	basic := feegrant.BasicAllowance{
		SpendLimit: sdktype.NewCoins(sdktype.NewCoin("ugas", sdktype.NewInt(400000000))),
	}

	grant = &basic

	msgGrant, err := feegrant.NewMsgGrantAllowance(grant, granter, grantee)
	if err != nil {
		//500
		log.Error("create account", "msg grant allowance error:", err.Error())
		return "", types.ErrInternal
	}

	baseTx := m.CreateBaseTx(pAddress, defultKeyPassword)
	res, err := m.BuildAndSend(sdktype.Msgs{msgGrant}, baseTx)
	if err != nil {
		//500
		log.Error("create account", "fee grant error:", err.Error())
		return "", types.ErrInternal
	}

	return res.Hash, nil
}

// EncodeData 加密序列
func (m Base) EncodeData(data string) string {
	hashBz := sha256.Sum256([]byte(data))
	hash := strings.ToUpper(hex.EncodeToString(hashBz[:]))
	return hash
}
