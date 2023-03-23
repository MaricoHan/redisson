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
	individualName := strings.TrimSpace(req.Individual.Name)
	certificateNum := strings.TrimSpace(req.Individual.CertificateNum)
	individualPhoneNum := strings.TrimSpace(req.Individual.PhoneNum)
	enterpriseName := strings.TrimSpace(req.Enterprise.Name)
	registrationNum := strings.TrimSpace(req.Enterprise.RegistrationNum)
	enterprisePhoneNum := strings.TrimSpace(req.Enterprise.PhoneNum)
	businessLicense := strings.TrimSpace(req.Enterprise.BusinessLicense)
	email := strings.TrimSpace(req.Enterprise.Email)

	authData := u.AuthData(ctx)
	params := dto.CreateUsers{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
		Usertype:   req.Usertype,
		Individual: &wallet.INDIVIDUALS{
			Name:            individualName,
			Region:          wallet.REGION(req.Individual.Region),
			CertificateType: wallet.CERTIFICATE_TYPE(req.Individual.CertificateType),
			CertificateNum:  certificateNum,
			PhoneNum:        individualPhoneNum,
		},
		Enterprise: &wallet.ENTERPRISES{
			Name:               enterpriseName,
			RegistrationRegion: wallet.REGION(req.Enterprise.RegistrationRegion),
			RegistrationNum:    registrationNum,
			PhoneNum:           enterprisePhoneNum,
			BusinessLicense:    businessLicense,
			Email:              email,
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
