package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
	"strings"

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

func (m Base) QueryRootAccount() (*models.TAccount, *types.AppError) {
	//platform address
	account, err := models.TAccounts(
		models.TAccountWhere.ProjectID.EQ(uint64(0)),
	).OneG(context.Background())
	if err != nil {
		//500
		log.Error("create account", "query root error:", err.Error())
		return nil, types.ErrInternal
	}
	return account, nil
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

func (m Base) CreateBaseTxSync(keyName, keyPassword string) sdktype.BaseTx {
	//from := "t_" + keyName
	return sdktype.BaseTx{
		From:     keyName,
		Gas:      m.gas,
		Fee:      m.coins,
		Mode:     sdktype.Sync,
		Password: keyPassword,
	}
}

func (m Base) BuildAndSign(msgs sdktype.Msgs, baseTx sdktype.BaseTx) ([]byte, string, error) {
	root, error := m.QueryRootAccount()
	if error != nil {
		return nil, "", error
	}
	baseTx.FeePayer = sdktype.AccAddress(root.Address)
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

// UndoTxIntoDataBase operationType : issue_class,mint_nft,edit_nft,edit_nft_batch,burn_nft,burn_nft_batch
func (m Base) UndoTxIntoDataBase(sender, operationType, taskId, txHash string, ProjectID uint64, signedData, message, tag []byte, gasUsed int64, exec boil.ContextExecutor) (uint64, error) {
	// Tx into database
	ttx := models.TTX{
		ProjectID:     ProjectID,
		Hash:          txHash,
		OriginData:    null.BytesFrom(signedData),
		OperationType: operationType,
		Status:        models.TTXSStatusUndo,
		Sender:        null.StringFrom(sender),
		Message:       null.JSONFrom(message),
		TaskID:        null.StringFrom(taskId),
		GasUsed:       null.Int64From(gasUsed),
		Tag:           null.JSONFrom(tag),
		Retry:         null.Int8From(0),
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
	l := len([]byte(class.ClassID)) + len([]byte(class.Status)) + len([]byte(class.Owner)) + len([]byte(class.TXHash)) + len([]byte(string(class.ProjectID))) + len([]byte(string(class.ID))) + len([]byte(string(class.Offset)))
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

func (m Base) Grant(address []string) (string, error) {
	root, error := m.QueryRootAccount()
	if error != nil {
		return "", error
	}
	granter, errs := sdktype.AccAddressFromBech32(root.Address)
	if errs != nil {
		//500
		log.Error("base account", "granter format error:", errs.Error())
		return "", types.ErrInternal
	}
	var msgs sdktype.Msgs
	for i := 0; i < len(address); i++ {
		grantee, errs := sdktype.AccAddressFromBech32(address[i])
		if errs != nil {
			//500
			log.Error("base account", "grantee format error:", errs.Error())
			return "", types.ErrInternal
		}
		var grant feegrant.FeeAllowanceI

		basic := feegrant.BasicAllowance{
			SpendLimit: nil,
		}

		grant = &basic

		msgGrant, err := feegrant.NewMsgGrantAllowance(grant, granter, grantee)
		if err != nil {
			//500
			log.Error("base account", "msg grant allowance error:", err.Error())
			return "", types.ErrInternal
		}
		msgs = append(msgs, msgGrant)
	}
	baseTx := m.CreateBaseTx(root.Address, defultKeyPassword)
	//动态计算gas
	baseTx.Gas = m.createAccount(int64(len(address)))
	res, err := m.BuildAndSend(msgs, baseTx)
	if err != nil {
		//500
		log.Error("base account", "fee grant error:", err.Error())
		return "", types.ErrInternal
	}
	return res.Hash, nil
}

// ValidateSigner validate signer
func (m Base) ValidateSigner(sender string, projectid uint64) error {
	//signer不能为project外账户
	_, err := models.TAccounts(
		models.TAccountWhere.ProjectID.EQ(projectid),
		models.TAccountWhere.Address.EQ(sender)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
		//404
		return types.ErrNotFound
	} else if err != nil {
		//500
		log.Error("validate signer", "query signer error:", err.Error())
		return types.ErrInternal
	}

	return nil
}

// ValidateRecipient validate recipient
func (m Base) ValidateRecipient(recipient string, projectid uint64) error {
	//recipient不能为project外的账户
	_, err := models.TAccounts(
		models.TAccountWhere.ProjectID.EQ(projectid),
		models.TAccountWhere.Address.EQ(recipient)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNoFound())) {
		//400
		return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientFound)
	} else if err != nil {
		//500
		log.Error("validate recipient", "query recipient error:", err.Error())
		return types.ErrInternal
	}
	return nil
}

// EncodeData 加密序列
func (m Base) EncodeData(data string) string {
	hashBz := sha256.Sum256([]byte(data))
	hash := strings.ToUpper(hex.EncodeToString(hashBz[:]))
	return hash
}

func (m Base) GasThan(address string, chainId, gas, platformId uint64) error {
	err := modext.Transaction(func(exec boil.ContextExecutor) error {
		tProjects, err := models.TProjects(
			models.TProjectWhere.PlatformID.EQ(null.Int64From(int64(platformId))),
		).All(context.Background(), exec)
		if err != nil {
			return errors.New("query PlatformID from TProjects failed")
		}
		var projects []uint64
		for _, v := range tProjects {
			projects = append(projects, v.ID)
		}
		// unPaidGas 待支付的gas
		tx, err := models.TTXS(
			qm.Select("SUM(gas_used) as gas_used"),
			models.TTXWhere.ProjectID.IN(projects),
			models.TTXWhere.Status.IN([]string{models.TTXSStatusPending, models.TTXSStatusUndo})).One(context.Background(), exec)
		if err != nil {
			return types.ErrNotFound
		}
		unPaidGas := tx.GasUsed.Int64
		chain, err := models.TChains(models.TChainWhere.ID.EQ(chainId),
			models.TChainWhere.Status.EQ(0)).OneG(context.Background())
		if err != nil {
			return types.ErrNotFound
		}
		//gasPrice 每条链的gasPrice
		gasPrice, ok := chain.GasPrice.Big.Float64()
		if !ok {
			return errors.New("cannot get float64 of gasPrice")
		}
		// unPaidMoney   = 这些未支付的交易需要扣除的money  =  gasPrice * unPaidGas
		unPaidMoney := float64(unPaidGas) * gasPrice
		pAccount, err := models.TPlatformAccounts(models.TPlatformAccountWhere.ID.EQ(platformId)).One(context.Background(), exec)

		if err != nil {
			return errors.New(fmt.Sprintf("cannot query PlatFormAccount and platformId is : %v", platformId))
		}
		//amount 平台方的余额
		amount, ok := pAccount.Amount.Big.Float64()
		if !ok {
			return errors.New("cannot get float64 of amount")
		}
		unPaidMoney = unPaidMoney + float64(gas)*gasPrice
		//如果amount小于未支付金额,返回错误
		if amount < unPaidMoney {
			return errors.New("balances not enough")
		}
		return err
	})
	return err
}
