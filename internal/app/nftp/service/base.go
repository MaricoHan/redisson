package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	ddc "github.com/bianjieai/ddc-sdk-go/app"
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

const (
	GasPrice         = 1e4
	GasLimit         = 1e6
	AuthorityAddress = "0x607F278304Fd91df7e2E6630a66809959c73978c"
	ChargeAddress    = "0xDdAEfC5E48a9ec1c63293997cea034570d5117c8"
	DDC721Address    = "0x87c263E5E1260eB02f9C5f7dE7504a91E324BBF0"
	DDC1155Address   = "0xf9E474ceD3486Bb003BE36cD1c41F4537b541c18"
)

var (
	SqlNotFound   = "records not exist"
	ClientBuilder = ddc.DDCSdkClientBuilder{}
	DDCClient     = ClientBuilder.
			SetGasPrice(GasPrice).
			SetGasLimit(GasLimit).
			SetAuthorityAddress(AuthorityAddress).
			SetChargeAddress(ChargeAddress).
			SetDDC721Address(DDC721Address).
			SetDDC1155Address(DDC1155Address).
			Build()
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

func (b Base) QueryRootAccount() (*models.TAccount, *types.AppError) {
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

func (b Base) CreateBaseTx(keyName, keyPassword string) sdktype.BaseTx {
	//from := "t_" + keyName
	return sdktype.BaseTx{
		From:     keyName,
		Gas:      b.gas,
		Fee:      b.coins,
		Mode:     sdktype.Commit,
		Password: keyPassword,
	}
}

func (b Base) CreateBaseTxSync(keyName, keyPassword string) sdktype.BaseTx {
	//from := "t_" + keyName
	return sdktype.BaseTx{
		From:     keyName,
		Gas:      b.gas,
		Fee:      b.coins,
		Mode:     sdktype.Sync,
		Password: keyPassword,
	}
}

func (b Base) BuildAndSign(msgs sdktype.Msgs, baseTx sdktype.BaseTx) ([]byte, string, error) {
	root, error := b.QueryRootAccount()
	if error != nil {
		return nil, "", error
	}
	baseTx.FeePayer = sdktype.AccAddress(root.Address)
	sigData, err := b.sdkClient.BuildAndSign(msgs, baseTx)
	if err != nil {
		return nil, "", err
	}
	hashBz := sha256.Sum256(sigData)
	hash := strings.ToUpper(hex.EncodeToString(hashBz[:]))
	return sigData, hash, nil
}

func (b Base) BuildAndSend(msgs sdktype.Msgs, baseTx sdktype.BaseTx) (sdktype.ResultTx, error) {
	sigData, err := b.sdkClient.BuildAndSend(msgs, baseTx)
	if err != nil {
		return sigData, err
	}
	return sigData, nil
}

// ValidateTx validate tx status
func (b Base) ValidateTx(hash string) (*models.TTX, error) {
	tx, err := models.TTXS(models.TTXWhere.Hash.EQ(hash)).OneG(context.Background())
	if err != nil {
		if err == sql.ErrNoRows || strings.Contains(err.Error(), SqlNotFound) {
			return tx, nil
		}
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

func (b Base) CreateGasMsg(inputAddress string, outputAddress []string) bank.MsgMultiSend {
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
func (b Base) MintNftsGas(originData []byte, amount uint64) uint64 {
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
func (b Base) CreateDenomGas(data []byte) uint64 {
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
func (b Base) TransferDenomGas(class *models.TClass) uint64 {
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
func (b Base) TransferOneNftGas(data []byte) uint64 {
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
func (b Base) TransferNftsGas(data []byte, amount uint64) uint64 {
	l := uint64(len(data))
	if l <= types.TransferMinNFTDataSize {
		return uint64(float64(types.TransferMinNFTGas) * config.Get().Chain.GasCoefficient)
	}
	res := (l-types.TransferMinNFTIncreaseDataSize*(amount-1)-types.TransferMinNFTDataSize)*types.TransferNFTCoefficient + types.TransferMinNFTGas + types.TransferMinNFTIncreaseGas*(amount-1)
	u := float64(res) * config.Get().Chain.GasCoefficient
	return uint64(u)
}

func (b Base) LenOfNft(tNft *models.TNFT) uint64 {
	len1 := len(tNft.Status + tNft.NFTID + tNft.Owner + tNft.ClassID + tNft.TXHash + tNft.Name.String + tNft.Metadata.String + tNft.URIHash.String + tNft.URI.String)
	len2 := 4 * 8 // 4 uint64
	len3 := len(tNft.CreateAt.String() + tNft.UpdateAt.String() + tNft.Timestamp.Time.String())
	return uint64(len1 + len2 + len3)
}

/**
Estimated gas required to edit nft
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049122
*/
func (b Base) EditNftGas(nftLen, signLen uint64) uint64 {
	gas := types.EditNFTBaseGas + types.EditNFTLenCoefficient*nftLen + types.EditNFTSignLenCoefficient*signLen
	res := float64(gas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to edit nfts
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049126
*/
func (b Base) EditBatchNftGas(nftLen, signLen uint64) uint64 {
	gas := types.EditBatchNFTBaseGas + types.EditBatchNFTLenCoefficient*nftLen + types.EditBatchNFTSignLenCoefficient*signLen
	res := float64(gas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to delete nft
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049119
*/
func (b Base) DeleteNftGas(nftLen uint64) uint64 {
	gas := types.DeleteNFTBaseGas + (nftLen-types.DeleteNFTBaseLen)*types.DeleteNFTCoefficient
	res := float64(gas) * config.Get().Chain.GasCoefficient
	return uint64(res)
}

/**
Estimated gas required to delete nfts
It is calculated as follows : http://wiki.bianjie.ai/pages/viewpage.action?pageId=58049124
*/
func (b Base) DeleteBatchNftGas(nftLen, n uint64) uint64 {
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
func (b Base) createAccount(count int64) uint64 {
	count -= 1
	res := types.CreateAccountGas + types.CreateAccountIncreaseGas*(count)
	u := float64(res) * config.Get().Chain.GasCoefficient
	return uint64(u)
}

func (b Base) Grant(address []string) (string, error) {
	root, error := b.QueryRootAccount()
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
	baseTx := b.CreateBaseTxSync(root.Address, config.Get().Server.DefaultKeyPassword)
	//动态计算gas
	baseTx.Gas = b.createAccount(int64(len(address)))
	res, err := b.BuildAndSend(msgs, baseTx)
	if err != nil {
		//500
		log.Error("base account", "fee grant error:", err.Error())
		return "", types.ErrInternal
	}
	return res.Hash, nil
}

// ValidateSigner validate signer
func (b Base) ValidateSigner(sender string, projectid uint64) error {
	//signer不能为project外账户
	_, err := models.TAccounts(
		models.TAccountWhere.ProjectID.EQ(projectid),
		models.TAccountWhere.Address.EQ(sender)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNotFound)) {
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
func (b Base) ValidateRecipient(recipient string, projectid uint64) error {
	//recipient不能为project外的账户
	_, err := models.TAccounts(
		models.TAccountWhere.ProjectID.EQ(projectid),
		models.TAccountWhere.Address.EQ(recipient)).OneG(context.Background())
	if (err != nil && errors.Cause(err) == sql.ErrNoRows) ||
		(err != nil && strings.Contains(err.Error(), SqlNotFound)) {
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
func (b Base) EncodeData(data string) string {
	hashBz := sha256.Sum256([]byte(data))
	hash := strings.ToUpper(hex.EncodeToString(hashBz[:]))
	return hash
}

func (b Base) GasThan(address string, chainId, gas, platformId uint64) error {
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
