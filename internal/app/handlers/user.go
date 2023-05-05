package handlers

import (
	"context"
	"strconv"
	"strings"

	"gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/wallet"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/utils/errors"
)

type IUser interface {
	CreateUsers(ctx context.Context, request interface{}) (interface{}, error)
	UpdateUsers(ctx context.Context, request interface{}) (interface{}, error)
	ShowUsers(ctx context.Context, request interface{}) (interface{}, error)
}

type User struct {
	Base
	PageBasic
	svc services.IUser
}

func NewUser(svc services.IUser) *User {
	return &User{svc: svc}
}

func (u User) CreateUsers(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*vo.CreateUserRequest)
	authData, err := u.validateUsers(ctx)
	if err != nil {
		return nil, err
	}
	// 校验参数
	name := strings.TrimSpace(req.Name)
	certificateNum := strings.TrimSpace(req.CertificateNum)
	phoneNum := strings.TrimSpace(req.PhoneNum)
	registrationNum := strings.TrimSpace(req.RegistrationNum)
	businessLicense := strings.TrimSpace(req.BusinessLicense)
	email := strings.TrimSpace(req.Email)

	authData.Code = constant.Wallet
	authData.Module = constant.Server

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
	authData, err := u.validateUsers(ctx)
	if err != nil {
		return nil, err
	}
	// 校验参数
	userId := strings.TrimSpace(req.UserId)
	phoneNum := strings.TrimSpace(req.PhoneNum)
	authData.Code = constant.Wallet
	authData.Module = constant.Server
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

func (u User) ShowUsers(ctx context.Context, request interface{}) (interface{}, error) {
	authData, err := u.validateUsers(ctx)
	if err != nil {
		return nil, err
	}
	authData.Code = constant.Wallet
	authData.Module = constant.Server

	params := dto.ShowUsers{
		ProjectID:  authData.ProjectId,
		ChainID:    authData.ChainId,
		Module:     authData.Module,
		AccessMode: authData.AccessMode,
		Code:       authData.Code,
		Usertype:   u.getUserType(ctx),
		UserCode:   strings.TrimSpace(u.getCode(ctx)),
	}

	return u.svc.ShowUsers(ctx, params)
}

func (u User) validateUsers(ctx context.Context) (vo.AuthData, error) {
	authData := u.AuthData(ctx)
	if !authData.ExistWalletService {
		return vo.AuthData{}, errors.ErrNotImplemented
	}
	return authData, nil
}

func (u User) getUserType(ctx context.Context) uint32 {
	userType := ctx.Value("user_type")
	if userType == nil || userType == "" {
		return 0
	}
	parseUint, err := strconv.ParseUint(userType.(string), 10, 64)
	if err != nil {
		return 0
	}
	return uint32(parseUint)
}

func (u User) getCode(ctx context.Context) string {
	code := ctx.Value("code")
	if code == nil || code == "" {
		return ""
	}
	return code.(string)
}
