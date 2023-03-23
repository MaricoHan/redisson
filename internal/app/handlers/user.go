package handlers

import (
	"context"
	"strings"

	"gitlab.bianjie.ai/avata/chains/api/pb/v2/wallet"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
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
	name := strings.TrimSpace(req.Name)
	certificateNum := strings.TrimSpace(req.CertificateNum)
	phoneNum := strings.TrimSpace(req.PhoneNum)
	registrationNum := strings.TrimSpace(req.RegistrationNum)
	businessLicense := strings.TrimSpace(req.BusinessLicense)
	email := strings.TrimSpace(req.Email)

	authData := u.AuthData(ctx)
	params := dto.CreateUsers{
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
		Usertype:   req.Usertype,
	}

	if req.Usertype == entity.UserTypeIndividual {
		params.Individual = &wallet.INDIVIDUALS{
			Name:            name,
			Region:          wallet.REGION(req.Region),
			CertificateType: wallet.CERTIFICATE_TYPE(req.CertificateType),
			CertificateNum:  certificateNum,
			PhoneNum:        phoneNum,
		}
	} else if req.Usertype == entity.UserTypeEnterprise {
		params.Enterprise = &wallet.ENTERPRISES{
			Name:               name,
			RegistrationRegion: wallet.REGION(req.RegistrationRegion),
			RegistrationNum:    registrationNum,
			PhoneNum:           phoneNum,
			BusinessLicense:    businessLicense,
			Email:              email,
		}
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
