package native

import (
	"context"
	"encoding/base64"
	"fmt"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"strings"

	sdkcrypto "github.com/irisnet/core-sdk-go/common/crypto"
	"github.com/irisnet/core-sdk-go/common/crypto/codec"
	"github.com/irisnet/core-sdk-go/common/crypto/hd"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
)

const algo = "secp256k1"
const hdPathPrefix = hd.BIP44Prefix + "0'/0/"

type nativeAccount struct {
	base *service.Base
}

func NewNFTAccount(base *service.Base) *service.AccountBase {
	return &service.AccountBase{
		Module:  service.NATIVE,
		Service: &nativeAccount{base: base},
	}
}

func (svc *nativeAccount) Create(params dto.CreateAccountP) (*dto.AccountRes, error) {
	// 写入数据库
	// sdk 创建账户
	var addresses []string

	err := modext.Transaction(func(exec boil.ContextExecutor) error {
		tAppOneObj, err := models.TApps(
			qm.SQL("SELECT * FROM `t_apps` WHERE (`t_apps`.`id` = ?) LIMIT 1 FOR UPDATE;", 1),
		).One(context.Background(), exec)
		if err != nil {
			//500
			log.Error("create account", "query app error:", err.Error())
			return types.ErrInternal
		}

		tAccounts := modext.TAccounts{}
		var i int64
		accOffsetStart := tAppOneObj.AccOffset
		for i = 0; i < params.Count; i++ {
			index := accOffsetStart + i
			hdPath := fmt.Sprintf("%s%d", hdPathPrefix, index)
			res, err := sdkcrypto.NewMnemonicKeyManagerWithHDPath(
				tAppOneObj.Mnemonic,
				config.Get().Chain.ChainEncryption,
				hdPath,
			)
			if err != nil {
				//500
				log.Debug("create account", "NewMnemonicKeyManagerWithHDPath error:", err.Error())
				return types.ErrInternal
			}
			_, priv := res.Generate()

			tmpAddress := sdktype.AccAddress(priv.PubKey().Address().Bytes()).String()

			tmp := &models.TAccount{
				ProjectID: params.ProjectID,
				Address:   tmpAddress,
				AccIndex:  uint64(index),
				PriKey:    base64.StdEncoding.EncodeToString(codec.MarshalPrivKey(priv)),
				PubKey:    base64.StdEncoding.EncodeToString(codec.MarshalPubkey(res.ExportPubKey())),
			}

			tAccounts = append(tAccounts, tmp)
			addresses = append(addresses, tmpAddress)
		}

		err = tAccounts.InsertAll(context.Background(), exec)
		if err != nil {
			log.Error("create account", "accounts insert error:", err.Error())
			return types.ErrInternal
		}
		tAppOneObj.AccOffset += params.Count
		updateRes, err := tAppOneObj.Update(context.Background(), exec, boil.Infer())
		if err != nil || updateRes == 0 {
			log.Error("create account", "apps insert error:", err.Error())
			return types.ErrInternal
		}
		// fee grant
		_, err = svc.base.Grant(addresses)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	result := &dto.AccountRes{}
	result.Accounts = addresses
	return result, nil
}

func (svc *nativeAccount) Show(params dto.AccountsP) (*dto.AccountsRes, error) {
	result := &dto.AccountsRes{}
	result.Offset = params.Offset
	result.Limit = params.Limit
	result.Accounts = []*dto.Account{}
	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TAccounts),
		qm.Select(models.TAccountColumns.Address, models.TAccountColumns.Gas),
		models.TAccountWhere.ID.NEQ(0),
		models.TAccountWhere.ProjectID.EQ(params.ProjectID),
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
		if strings.Contains(err.Error(), service.SqlNotFound) {
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

func (svc *nativeAccount) History(params dto.AccountsP) (*dto.AccountOperationRecordRes, error) {
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
		models.TMSGWhere.ProjectID.EQ(params.ProjectID),
		models.TMSGWhere.Operation.NEQ(models.TMSGSOperationSysIssueClass),
	}

	if params.Account != "" {
		queryMod = append(queryMod, qm.Where("signer = ? OR recipient = ?", params.Account, params.Account))
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
		if strings.Contains(err.Error(), service.SqlNotFound) {
			return result, nil
		}
		log.Error("account history", "query error:", err)
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
