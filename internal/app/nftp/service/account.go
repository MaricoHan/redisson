package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"strings"
	"time"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/irisnet/core-sdk-go/common/crypto/hd"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"

	"github.com/irisnet/core-sdk-go/common/crypto/codec"
	"github.com/volatiletech/null/v8"

	"github.com/irisnet/core-sdk-go/bank"
	sdkcrypto "github.com/irisnet/core-sdk-go/common/crypto"
	sdktype "github.com/irisnet/core-sdk-go/types"
)

const algo = "secp256k1"
const hdPathPrefix = hd.BIP44Prefix + "0'/0/"

const defultKeyPassword = "12345678"

type Account struct {
	base *Base
}

func NewAccount(base *Base) *Account {
	return &Account{base: base}
}

func (svc *Account) CreateAccount(params dto.CreateAccountP) ([]string, error) {
	// 写入数据库
	// sdk 创建账户
	var addresses []string
	classOne, err := models.TAccounts(
		models.TAccountWhere.AppID.EQ(uint64(0)),
	).OneG(context.Background())
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		//400
		return nil, types.NewAppError(types.RootCodeSpace, types.QueryDataFailed, "root account not found")
	} else if err != nil {
		//500
		log.Error("create account", "query root account error:", err.Error())
		return nil, types.ErrCreate
	}
	tmsgs := modext.TMSGs{}
	var msgs bank.MsgMultiSend
	var resultTx sdktype.ResultTx
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		tAppOneObj, err := models.TApps(models.TAppWhere.ID.EQ(params.AppID)).One(context.Background(), exec)
		if err != nil && errors.Cause(err) == sql.ErrNoRows {
			//400
			return types.NewAppError(types.RootCodeSpace, types.QueryDataFailed, "app not found")
		} else if err != nil {
			//500
			log.Error("create account", "query app error:", err.Error())
			return types.ErrCreate
		}

		tAccounts := modext.TAccounts{}
		var i int64
		accOffsetStart := tAppOneObj.AccOffset
		for i = 0; i < params.Count; i++ {
			index := accOffsetStart + i
			hdPath := fmt.Sprintf("%s%d", hdPathPrefix, index)
			res, err := sdkcrypto.NewMnemonicKeyManagerWithHDPath(
				tAppOneObj.Mnemonic,
				algo,
				hdPath,
			)
			if err != nil {
				//500
				log.Debug("create account", "NewMnemonicKeyManagerWithHDPath error:", err.Error())
				return types.ErrCreate
			}
			_, priv := res.Generate()

			//privStr, err := res.ExportPrivKey(keyPassword)
			//if err != nil {
			//	return nil, types.ErrAccountCreate
			//}

			tmpAddress := sdktype.AccAddress(priv.PubKey().Address().Bytes()).String()

			tmp := &models.TAccount{
				AppID:    params.AppID,
				Address:  tmpAddress,
				AccIndex: uint64(index),
				PriKey:   base64.StdEncoding.EncodeToString(codec.MarshalPrivKey(priv)),
				PubKey:   base64.StdEncoding.EncodeToString(codec.MarshalPubkey(res.ExportPubKey())),
			}

			tAccounts = append(tAccounts, tmp)
			addresses = append(addresses, tmpAddress)
		}
		err = tAccounts.InsertAll(context.Background(), exec)
		if err != nil {
			log.Debug("create account", "accounts insert error:", err.Error())
			return types.ErrCreate
		}
		tAppOneObj.AccOffset += params.Count
		updateRes, err := tAppOneObj.Update(context.Background(), exec, boil.Infer())
		if err != nil || updateRes == 0 {
			return types.ErrInternal
		}
		msgs = svc.base.CreateGasMsg(classOne.Address, addresses)
		tx := svc.base.CreateBaseTx(classOne.Address, defultKeyPassword)
		resultTx, err = svc.base.BuildAndSend(sdktype.Msgs{&msgs}, tx)
		if err != nil {
			log.Error("create account", "build and send, error:", err)
			return types.ErrBuildAndSend
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range msgs.Outputs {
		message := map[string]string{
			"recipient": v.Address,
			"amount":    v.Coins.String()[0 : len(v.Coins.String())-6],
		}
		messageByte, _ := json.Marshal(message)
		tmsgs = append(tmsgs, &models.TMSG{
			AppID:     params.AppID,
			TXHash:    resultTx.Hash,
			Timestamp: null.TimeFrom(time.Now()),
			Module:    "account",
			Operation: "add_gas",
			Signer:    classOne.Address,
			Recipient: null.StringFrom(v.Address),
			Message:   messageByte,
		})
	}
	err = tmsgs.InsertAll(context.Background(), boil.GetContextDB())
	if err != nil {
		log.Error("create account", "msgs create error:", err)
		return nil, types.ErrCreate
	}
	return addresses, nil
}

func (svc *Account) Accounts(params dto.AccountsP) (*dto.AccountsRes, error) {
	result := &dto.AccountsRes{}
	result.Offset = params.Offset
	result.Limit = params.Limit
	result.Accounts = []*dto.Account{}
	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TAccounts),
		qm.Select(models.TAccountColumns.Address, models.TAccountColumns.Gas),
		models.TAccountWhere.AppID.EQ(params.AppID),
	}
	if params.Account != "" {
		queryMod = append(queryMod, models.TAccountWhere.Address.EQ(params.Account))
	}

	if params.StartDate != nil {
		queryMod = append(queryMod, models.TAccountWhere.CreateAt.GTE(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TAccountWhere.CreateAt.LTE(*params.EndDate))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TAccountColumns.CreateAt)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TAccountColumns.CreateAt)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}

	var modelResults []*models.TAccount
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
		if strings.Contains(err.Error(), "records not exist") {
			return result, nil
		}
		return nil, types.ErrInternal
	}

	result.TotalCount = total
	var accounts []*dto.Account
	for _, modelResult := range modelResults {
		account := &dto.Account{
			Account: modelResult.Address,
			Gas:     modelResult.Gas.Uint64,
		}
		accounts = append(accounts, account)
	}
	result.Accounts = accounts

	return result, nil
}

func (svc *Account) AccountsHistory(params dto.AccountsP) (*dto.AccountOperationRecordRes, error) {
	result := &dto.AccountOperationRecordRes{
		PageRes: dto.PageRes{
			Offset:     params.Offset,
			Limit:      params.Limit,
			TotalCount: 0,
		},
		OperationRecords: []*dto.AccountOperationRecords{},
	}
	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TMSGS),
		models.TMSGWhere.AppID.EQ(params.AppID),
	}

	if params.Account != "" {
		queryMod = append(queryMod, models.TMSGWhere.Signer.EQ(params.Account))
	}
	if params.Module != "" {
		queryMod = append(queryMod, models.TMSGWhere.Module.EQ(params.Module))
	}
	if params.Operation != "" {
		queryMod = append(queryMod, models.TMSGWhere.Operation.EQ(params.Operation))
	}
	if params.StartDate != nil {
		queryMod = append(queryMod, models.TMSGWhere.CreateAt.GTE(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TMSGWhere.CreateAt.LTE(*params.EndDate))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TMSGColumns.Timestamp)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TMSGColumns.Timestamp)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}

	var modelResults []*models.TMSG
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
		if strings.Contains(err.Error(), "records not exist") {
			return result, nil
		}
		return nil, types.ErrInternal
	}

	result.TotalCount = total
	var accountOperationRecords []*dto.AccountOperationRecords
	for _, modelResult := range modelResults {
		accountOperationRecord := &dto.AccountOperationRecords{
			TxHash:    modelResult.TXHash,
			Module:    modelResult.Module,
			Operation: modelResult.Operation,
			Signer:    modelResult.Signer,
			Timestamp: modelResult.Timestamp.Time.String(),
			Message:   modelResult.Message,
		}
		accountOperationRecords = append(accountOperationRecords, accountOperationRecord)
	}

	result.OperationRecords = accountOperationRecords
	return result, nil
}
