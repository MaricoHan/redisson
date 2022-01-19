package handlers

import (
	"context"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
)

type IAccount interface {
	CreateAccount(ctx context.Context, _ interface{}) (interface{}, error)
	Accounts(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewAccount() IAccount {
	return newAccount()
}

type account struct {
	base
	svc *service.Account
}

func newAccount() *account {
	return &account{}
}

// CreateAccount Create one or more accounts
// return creation result
func (h account) CreateAccount(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(vo.CreateAccountRequest)
	params := dto.CreateAccountP{
		AppID: h.AppID(ctx),
		Count: req.Count,
	}
	// 校验参数 end
	return h.svc.CreateAccount(params)
}

// Accounts return account list
func (h account) Accounts(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	params := dto.AccountsP{
		AppID: h.AppID(ctx),
		// todo
	}
	// 校验参数 end

	// 业务数据入库的地方
	return h.svc.Accounts(params)
}
