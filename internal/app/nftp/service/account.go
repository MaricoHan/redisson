package service

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/irisnet/core-sdk-go/common/crypto/hd"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/modext"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

	sdkcrypto "github.com/irisnet/core-sdk-go/common/crypto"
)

const algo = "secp256k1"
const hdPathPrefix = hd.BIP44Prefix + "0'/0/"

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
			return nil, types.ErrAccountCreate
		}
		_, prv := res.Generate()
		tmpAddress := prv.PubKey().Address().String()
		tmp := &models.TAccount{
			AppID:   params.AppID,
			Address: tmpAddress,
			PriKey:  string(prv.Bytes()),
			PubKey:  string(prv.PubKey().Bytes()),
		}

		tAccounts = append(tAccounts, tmp)
		addresses = append(addresses, tmpAddress)
	}
	err = tAccounts.InsertAll(context.Background(), db)
	if err != nil {
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
	queryMod := []qm.QueryMod{
		qm.Select(models.TAccountColumns.Address, models.TAccountColumns.Gas),
		models.TAccountWhere.AppID.EQ(params.AppID),
	}
	if params.Account != "" {
		queryMod = append(queryMod, models.TAccountWhere.Address.EQ(params.Account))
	}

	if params.StartDate != nil {
		queryMod = append(queryMod, models.TAccountWhere.CreateAt.EQ(*params.StartDate))
	}
	if params.EndDate != nil {
		queryMod = append(queryMod, models.TAccountWhere.CreateAt.EQ(*params.EndDate))
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
	total, err := modext.PageQuery(
		context.Background(),
		orm.GetDB(),
		queryMod,
		&modelResults,
		params.Offset,
		params.Limit,
	)
	if err != nil {
		return nil, types.ErrMysqlConn
	}
	result := &dto.AccountsRes{}
	result.Offset = params.Offset
	result.Limit = params.Limit
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
