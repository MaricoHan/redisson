package wenchangchain_ddc

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethereumcrypto "github.com/ethereum/go-ethereum/crypto"
	sdkcrypto "github.com/irisnet/core-sdk-go/common/crypto"
	"github.com/irisnet/core-sdk-go/common/crypto/codec"
	ethsecp256k1 "github.com/irisnet/core-sdk-go/common/crypto/keys/eth_secp256k1"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"strings"

	"github.com/irisnet/core-sdk-go/common/crypto/hd"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"
)

type ddcAccount struct {
	base *service.Base
}

func NewDDCAccount(base *service.Base) *service.AccountBase {
	return &service.AccountBase{
		Module: service.DDC,
		Service: &ddcAccount{
			base: base,
		},
	}
}

const hdPathPrefix = hd.BIP44Prefix + "0'/0/"

func (svc *ddcAccount) Create(params dto.CreateAccountP) (*dto.AccountRes, error) {
	// 写入数据库
	// sdk 创建账户
	var addresses []string
	var bech32addresses []string
	client:=service.NewDDCClient()
	err := modext.Transaction(func(exec boil.ContextExecutor) error {
		tAppOneObj, err := models.TConfigs(
			qm.SQL("SELECT * FROM `t_configs` WHERE (`t_configs`.`id` = ?) LIMIT 1 FOR UPDATE;", 1),
		).One(context.Background(), exec)
		if err != nil {
			//500
			log.Error("create ddc account", "query app error:", err.Error())
			return types.ErrInternal
		}

		tAccounts := modext.TDDCAccounts{}
		var i int64
		accOffsetStart := tAppOneObj.AccOffset
		mnemonic64, err := base64.StdEncoding.DecodeString(tAppOneObj.Mnemonic)
		if err != nil {
			log.Error("create account", "mnemonic base64 error:", err.Error())
			return types.ErrInternal
		}
		mnemonic, err := types.Decrypt(mnemonic64, config.Get().Server.DefaultKeyPassword)
		if err != nil {
			log.Error("create account", "mnemonic Decrypt error:", err.Error())
			return types.ErrInternal
		}

		//查询有授权权限账户
		owner, err := models.TDDCAccounts(
			models.TDDCAccountWhere.ProjectID.EQ(uint64(0)),
		).OneG(context.Background())
		if err != nil {
			//500
			log.Error("create account", "query owner error:", err.Error())
			return types.ErrInternal
		}

		for i = 0; i < params.Count; i++ {

			index := accOffsetStart +i
			hdPath := fmt.Sprintf("%s%d", hdPathPrefix, index)
			res, err := sdkcrypto.NewMnemonicKeyManagerWithHDPath(
				mnemonic,
				config.Get().Chain.DdcEncryption,
				hdPath,
			)
			if err != nil {
				//500
				log.Debug("create ddc account", "NewMnemonicKeyManagerWithHDPath error:", err.Error())
				return types.ErrInternal
			}
			_, priv := res.Generate()
			tmpAddress := sdktype.AccAddress(priv.PubKey().Address().Bytes()).String()

			//Converts key to Ethermint secp256k1 implementation
			ethPrivKey, ok := priv.(*ethsecp256k1.PrivKey)
			if !ok {
				return fmt.Errorf("invalid private key type %T, expected %T", priv,&ethsecp256k1.PrivKey{})
			}
			keys, err := ethPrivKey.ToECDSA()
			if err != nil {
				return err
			}

			// Formats key for output
			privB := ethereumcrypto.FromECDSA(keys)
			keyS := strings.ToUpper(hexutil.Encode(privB)[2:])

			//私钥加密
			priKey, err := types.Encrypt(keyS, config.Get().Server.DefaultKeyPassword)
			if err != nil {
				log.Error("create account", "encrypt prikey error:", err.Error())
				return types.ErrInternal
			}

			//hex address
			ddc721 := client.GetDDC721Service(true)
			addr, err := ddc721.Bech32ToHex(tmpAddress)
			if err != nil {
				return err
			}

			// add did
			authority := client.GetAuthorityService()
			_, err = authority.AddAccountByPlatform(owner.Address, addr, addr, "did:"+addr)
			if err != nil {
				return err
			}

			//recharge
			charge:=client.GetChargeService()
			_, err =charge.Recharge(owner.Address,addr,20)
			if err != nil {
				return err
			}

			tmp := &models.TDDCAccount{
				ProjectID: params.ProjectID,
				Address:   addr,
				AccIndex:  uint64(index),
				PriKey:    base64.StdEncoding.EncodeToString(priKey),
				PubKey:    base64.StdEncoding.EncodeToString(codec.MarshalPubkey(res.ExportPubKey())),
				Did: null.StringFrom("did:"+addr),
			}

			tAccounts = append(tAccounts, tmp)
			addresses = append(addresses, addr)
			bech32addresses = append(bech32addresses,tmpAddress)
		}

		//send balance
		root, err := svc.base.QueryRootAccount()
		if err != nil {
			return err
		}
		msgs := svc.base.CreateGasMsg(root.Address, bech32addresses)
		tx := svc.base.CreateBaseTxSync(root.Address, "")
		tx.Gas = svc.base.CreateAccount(params.Count)
		_, err = svc.base.BuildAndSend(sdktype.Msgs{&msgs}, tx)
		if err != nil {
			log.Error("create account", "build and send, error:", err)
			return types.ErrBuildAndSend
		}

		err = tAccounts.InsertAll(context.Background(), exec)
		if err != nil {
			log.Error("create ddc account", "accounts insert error:", err.Error())
			return types.ErrInternal
		}
		tAppOneObj.AccOffset += params.Count
		updateRes, err := tAppOneObj.Update(context.Background(), exec, boil.Infer())
		if err != nil || updateRes == 0 {
			log.Error("create ddc account", "apps insert error:", err.Error())
			return types.ErrInternal
		}
		// fee grant
		_, err = svc.base.Grant(bech32addresses)
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

func (svc *ddcAccount) Show(params dto.AccountsP) (*dto.AccountsRes, error) {
	result := &dto.AccountsRes{}
	result.Offset = params.Offset
	result.Limit = params.Limit
	result.Accounts = []*dto.Account{}
	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TDDCAccounts),
		qm.Select(models.TDDCAccountColumns.Address, models.TDDCAccountColumns.Gas),
		models.TDDCAccountWhere.ID.NEQ(0),
		models.TDDCAccountWhere.ProjectID.EQ(params.ProjectID),
	}
	if params.Account != "" {
		queryMod = append(queryMod, models.TDDCAccountWhere.Address.EQ(params.Account))
	}

	if params.StartDate != nil {
		queryMod = append(queryMod, models.TDDCAccountWhere.CreateAt.GTE(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TDDCAccountWhere.CreateAt.LTE(*params.EndDate))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TDDCAccountColumns.CreateAt)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TDDCAccountColumns.CreateAt)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}

	var modelResults []*models.TDDCAccount
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
		log.Error("account show", "query error:", err)
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

func (svc *ddcAccount) History(params dto.AccountsP) (*dto.AccountOperationRecordRes, error) {
	result := &dto.AccountOperationRecordRes{
		PageRes: dto.PageRes{
			Offset:     params.Offset,
			Limit:      params.Limit,
			TotalCount: 0,
		},
		OperationRecords: []*dto.AccountOperationRecords{},
	}
	queryMod := []qm.QueryMod{
		qm.From(models.TableNames.TDDCMSGS),
		models.TDDCMSGWhere.ProjectID.EQ(params.ProjectID),
		models.TDDCMSGWhere.Operation.NEQ(models.TDDCMSGSOperationSysIssueClass),
	}

	if params.Account != "" {
		queryMod = append(queryMod, qm.Where("signer = ? OR recipient = ?", params.Account, params.Account))
	}
	if params.Module != "" {
		queryMod = append(queryMod, models.TDDCMSGWhere.Module.EQ(params.Module))
	}
	if params.Operation != "" {
		queryMod = append(queryMod, models.TDDCMSGWhere.Operation.EQ(params.Operation))
	}
	if params.StartDate != nil {
		queryMod = append(queryMod, models.TDDCMSGWhere.CreateAt.GTE(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TDDCMSGWhere.CreateAt.LTE(*params.EndDate))
	}
	if params.SortBy != "" {
		orderBy := ""
		switch params.SortBy {
		case "DATE_DESC":
			orderBy = fmt.Sprintf("%s DESC", models.TDDCMSGColumns.Timestamp)
		case "DATE_ASC":
			orderBy = fmt.Sprintf("%s ASC", models.TDDCMSGColumns.Timestamp)
		}
		queryMod = append(queryMod, qm.OrderBy(orderBy))
	}

	var modelResults []*models.TDDCMSG
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
