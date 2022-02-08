package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/irisnet/core-sdk-go/common/crypto/hd"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"gitlab.bianjie.ai/irita-paas/open-api/config"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"

	"github.com/irisnet/core-sdk-go/common/crypto/codec"
	"github.com/volatiletech/null/v8"
	"golang.org/x/sync/errgroup"

	"github.com/irisnet/core-sdk-go/bank"
	sdkcrypto "github.com/irisnet/core-sdk-go/common/crypto"
	sdktype "github.com/irisnet/core-sdk-go/types"
	sqltype "github.com/volatiletech/sqlboiler/v4/types"
	http2 "gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/http"
)

const algo = "secp256k1"
const hdPathPrefix = hd.BIP44Prefix + "0'/0/"

const defultKeyPassword = "12345678"

type BsnAccount struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
}

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
	if err != nil {
		return nil, types.ErrNotFound
	}
	tmsgs := modext.TMSGs{}
	var msgs bank.MsgMultiSend
	var resultTx sdktype.ResultTx
	env := config.Get().Server.Env
	err = modext.Transaction(func(exec boil.ContextExecutor) error {
		tAppOneObj, err := models.TApps(models.TAppWhere.ID.EQ(params.AppID)).One(context.Background(), exec)
		if err != nil {
			return types.ErrNotFound
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
				log.Debug("create account", "NewMnemonicKeyManagerWithHDPath error:", err.Error())
				return types.ErrAccountCreate
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
			return types.ErrAccountCreate
		}
		tAppOneObj.AccOffset += params.Count
		updateRes, err := tAppOneObj.Update(context.Background(), exec, boil.Infer())
		if err != nil || updateRes == 0 {
			return types.ErrInternal
		}
		// create chain account
		if env == "stage" {
			msgs = svc.base.CreateGasMsg(classOne.Address, addresses)
			tx := svc.base.CreateBaseTx(classOne.Address, defultKeyPassword)
			resultTx, err = svc.base.BuildAndSend(sdktype.Msgs{&msgs}, tx)
			if err != nil {
				log.Error("create account", "build and send, error:", err)
				return types.ErrBuildAndSend
			}
		} else {
			group := new(errgroup.Group)
			for _, v := range tAccounts {
				group.Go(func() error {
					var bsnAccount BsnAccount
					chainClient := map[string]interface{}{
						"chainClientName": fmt.Sprintf("%s%d%d", tAppOneObj.Name.String, tAppOneObj.ID, v.AccIndex),
						"chainClientAddr": v.Address,
					}
					url := fmt.Sprintf("%s%s", config.Get().Server.BSNUrl, fmt.Sprintf("/api/%s/account/generate", config.Get().Server.BSNProjectId))
					res, err := http2.Post(url, "application/json", chainClient)
					if err != nil {
						return err
					}
					defer res.Body.Close()
					body, err := ioutil.ReadAll(res.Body)
					json.Unmarshal(body, &bsnAccount)
					if bsnAccount.Code != 0 {
						return errors.New(bsnAccount.Message)
					}
					return nil
				})
			}
			if err := group.Wait(); err != nil {
				log.Error("create account", "group, error:", err)
				return types.ErrAccountCreate
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if env == "stage" {
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
				Message:   sqltype.JSON(messageByte),
			})
		}
		err = tmsgs.InsertAll(context.Background(), boil.GetContextDB())
		if err != nil {
			log.Error("create account", "msgs create error:", err)
			return nil, types.ErrAccountCreate
		}
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
