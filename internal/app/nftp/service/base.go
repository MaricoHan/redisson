package service

import (
	sdk "github.com/irisnet/core-sdk-go"
	sdktype "github.com/irisnet/core-sdk-go/types"
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

func (m Base) CreateBaseTx(keyName, keyPassword string) sdktype.BaseTx {
	return sdktype.BaseTx{
		From:     keyName,
		Gas:      m.gas,
		Fee:      m.coins,
		Mode:     sdktype.Async,
		Password: keyPassword,
	}
}

func (m Base) BuildAndSign(msgs sdktype.Msgs, baseTx sdktype.BaseTx) ([]byte, error) {
	return m.sdkClient.BuildAndSign(msgs, baseTx)
}
