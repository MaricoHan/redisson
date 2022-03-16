package service

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"strings"
)

type SignListener struct {
}

// SignEvent 用户自定义的签名方法
func (s *SignListener) SignEvent(sender common.Address, tx *types.Transaction) (*types.Transaction, error) {
	// 提取私钥
	privateKey, err := StringToPrivateKey("0x" + strings.ToUpper(sender.Hex()[2:]))
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
