package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/irisnet/core-sdk-go/common/crypto/codec"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/irisnet/core-sdk-go/common/crypto/hd"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"

	sdkcrypto "github.com/irisnet/core-sdk-go/common/crypto"
	sdktype "github.com/irisnet/core-sdk-go/types"
)

const algo = "secp256k1"
const hdPathPrefix = hd.BIP44Prefix + "0'/0/"

const defultKeyPassword = "12345678"

type Account struct {
}

func NewAccount() *Account {
	return &Account{}
}

func (svc *Account) CreateAccount(params dto.CreateAccountP) ([]string, error) {
	// 写入数据库
	// sdk 创建账户
	db, err := orm.GetDB().Begin()
	if err != nil {
		return nil, types.ErrMysqlConn
	}
	tAppOneObj, err := models.TApps(models.TAppWhere.ID.EQ(params.AppID)).One(context.Background(), db)
	if err != nil {
		return nil, types.ErrInternal
	}

	tAccounts := modext.TAccounts{}

	var addresses []string
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
			log.Debug("create account", "NewMnemonicKeyManagerWithHDPath error:", err.Error())
			return nil, types.ErrAccountCreate
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
	err = tAccounts.InsertAll(context.Background(), db)
	if err != nil {
		log.Debug("create account", "accounts insert error:", err.Error())
		return nil, types.ErrAccountCreate
	}

	tAppOneObj.AccOffset += params.Count
	updateRes, err := tAppOneObj.Update(context.Background(), db, boil.Infer())
	if err != nil || updateRes == 0 {
		return nil, types.ErrInternal
	}
	err = db.Commit()
	if err != nil {
		return nil, types.ErrInternal
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
			orderBy = fmt.Sprintf("%s desc", models.TAccountColumns.CreateAt)
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

		return nil, types.ErrMysqlConn
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
			orderBy = fmt.Sprintf("%s desc", models.TMSGColumns.Timestamp)
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

		return nil, types.ErrMysqlConn
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
