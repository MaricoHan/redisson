package chain

import (
	"context"
	"database/sql"
	"encoding/base64"
	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"

	keystore "github.com/irisnet/core-sdk-go/types/store"
)

const algo = "secp256k1"

type MysqlKeyDao struct {
	db   *sql.DB
	algo string
}

func NewMsqlKeyDao(db *sql.DB, algo string) MysqlKeyDao {
	return MysqlKeyDao{db: db, algo: algo}
}

// Write will use user password to encrypt data and save to file, the file name is user name
func (k MysqlKeyDao) Write(name, password string, store keystore.KeyInfo) error {
	panic("not yet implemented")
}

// Read will read encrypted data from file and decrypt with user password
func (k MysqlKeyDao) Read(name, password string) (keystore.KeyInfo, error) {
	tAccountOneObj, err := models.TAccounts(
		qm.Select(
			models.TAccountColumns.Address,
			models.TAccountColumns.PriKey,
			models.TAccountColumns.PubKey,
		),
		models.TAccountWhere.Address.EQ(name),
	).One(context.Background(), k.db)
	if err != nil {
		return keystore.KeyInfo{}, err
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(tAccountOneObj.PubKey)
	if err != nil {
		return keystore.KeyInfo{}, err
	}

	priKeyDecodeString, err := base64.StdEncoding.DecodeString(tAccountOneObj.PriKey)
	if err != nil {
		return keystore.KeyInfo{}, err
	}
	priKeyStrings, err := types.Decrypt(priKeyDecodeString, config.Get().Server.DefaultKeyPassword)
	if err != nil {
		return keystore.KeyInfo{}, err
	}
	priKey, err := base64.StdEncoding.DecodeString(priKeyStrings)
	if err != nil {
		return keystore.KeyInfo{}, err
	}
	store := keystore.KeyInfo{
		Name:         name,
		Algo:         k.algo,
		PrivKeyArmor: string(priKey),
		PubKey:       pubKeyBytes,
	}

	return store, nil
}

// Delete will delete user data and use user password to verify permissions
func (k MysqlKeyDao) Delete(name, password string) error {
	panic("not yet implemented")
}

// Has returns whether the specified user name exists
func (k MysqlKeyDao) Has(name string) bool {
	exists, err := models.TAccounts(models.TAccountWhere.Address.EQ(name)).Exists(context.Background(), k.db)
	if err != nil {
		return false
	}
	return exists
}
