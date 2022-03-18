package service

import (
	"log"
	"strings"

	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	types2 "gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"gitlab.bianjie.ai/irita-paas/open-api/config"
	log2 "gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

type SignListener struct {
}

// SignEvent 用户自定义的签名方法
func (s *SignListener) SignEvent(sender common.Address, tx *types.Transaction) (*types.Transaction, error) {
	account, err := models.TDDCAccounts(
		models.TDDCAccountWhere.Address.EQ("0x"+strings.ToUpper(sender.Hex()[2:])),
	).OneG(context.Background())
	if err != nil {
		return nil, types.ErrInvalidSig
	}
	priKey, err := base64.StdEncoding.DecodeString(account.PriKey)
	if err != nil {
		log2.Error("sign event", "priKey base64 error:", err.Error())
		return nil,types2.ErrInternal
	}
	prKey, err := types2.Decrypt(priKey, config.Get().Server.DefaultKeyPassword)
	if err != nil {
		log2.Error("sign event", "priKey Decrypt error:", err.Error())
		return nil,types2.ErrInternal
	}
	//提取私钥
	privateKey, err := StringToPrivateKey("0x"+prKey)
	if err != nil {
		log.Fatalf("StringToPrivateKey failed:%v", err)
	}
	// 签名
	signTx, err := types.SignTx(tx, &types.HomesteadSigner{}, privateKey)
	return signTx, err
}

// StringToPrivateKey 从明文的私钥字符串转换成该类型
func StringToPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyByte, err := hexutil.Decode(privateKeyStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
