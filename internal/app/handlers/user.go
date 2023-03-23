package handlers

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/chains/api/pb/v2/wallet"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
)

type IUser interface {
	CreateUsers(ctx context.Context, request interface{}) (interface{}, error)
	UpdateUsers(ctx context.Context, request interface{}) (interface{}, error)
}

type User struct {
	base
	pageBasic
	svc services.IUser
}

func NewUser(svc services.IUser) *User {
	return &User{svc: svc}
}

func (u User) CreateUsers(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.CreateUserRequest)

	// 校验参数
	req.Individual.Name = strings.TrimSpace(req.Individual.Name)

	authData := u.AuthData(ctx)
	params := dto.CreateUsers{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
		Usertype:   req.Usertype,
		Individual: &wallet.INDIVIDUALS{
			Name:            req.Individual.Name,
			Region:          wallet.REGION(req.Individual.Region),
			CertificateType: wallet.CERTIFICATE_TYPE(req.Individual.CertificateType),
			CertificateNum:  req.Individual.CertificateNum,
			PhoneNum:        req.Individual.PhoneNum,
		},
		Enterprise: &wallet.ENTERPRISES{
			Name:               req.Enterprise.Name,
			RegistrationRegion: wallet.REGION(req.Enterprise.RegistrationRegion),
			RegistrationNum:    req.Enterprise.RegistrationNum,
			PhoneNum:           req.Enterprise.PhoneNum,
			BusinessLicense:    req.Enterprise.BusinessLicense,
			Email:              req.Enterprise.Email,
		},
	}
	return u.svc.CreateUsers(ctx, params)
}

func (u User) UpdateUsers(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.UpdateUserRequest)

	// 校验参数
	userId := strings.TrimSpace(req.UserId)
	phoneNum := strings.TrimSpace(req.PhoneNum)
	authData := u.AuthData(ctx)
	params := dto.UpdateUsers{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
		UserId:     userId,
		PhoneNum:   phoneNum,
	}
	return u.svc.UpdateUsers(ctx, params)
}
