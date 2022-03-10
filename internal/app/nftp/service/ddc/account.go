package ddc

import (
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
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

func (svc *ddcAccount) Create(params dto.CreateAccountP) (*dto.AccountRes, error) {
	return nil, nil
}

func (svc *ddcAccount) Show(params dto.AccountsP) (*dto.AccountsRes, error) {
	return nil, nil
}

func (svc *ddcAccount) History(params dto.AccountsP) (*dto.AccountOperationRecordRes, error) {
	panic("implement me")
}
