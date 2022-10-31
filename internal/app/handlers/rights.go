package handlers

import (
	"context"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	service "gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"strings"
)

type IRights interface {
	Register(ctx context.Context, request interface{}) (response interface{}, err error)
	EditRegister(ctx context.Context, request interface{}) (response interface{}, err error)
	QueryRegister(ctx context.Context, request interface{}) (response interface{}, err error)
	UserAuth(ctx context.Context, request interface{}) (response interface{}, err error)
	EditUserAuth(ctx context.Context, request interface{}) (response interface{}, err error)
	QueryUserAuth(ctx context.Context, request interface{}) (response interface{}, err error)

	Dict(ctx context.Context, request interface{}) (response interface{}, err error)
	Region(ctx context.Context, request interface{}) (response interface{}, err error)
}

type Rights struct {
	base
	svc service.IRights
}

func NewRights(svc service.IRights) *Rights {
	return &Rights{svc: svc}
}

func (r Rights) Register(ctx context.Context, request interface{}) (response interface{}, err error) {
	req := request.(*vo.RegisterRequest)

	// 校验参数
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	// todo 仿照tag处理metadata
	//tagBytes, err := r.ValidateTag(req.Metadata)
	//if err != nil {
	//	return nil, err
	//}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	authData := r.AuthData(ctx)
	var authorsIndividuals []dto.Individual
	var authorsCorporates []dto.Corporate
	var copyRightsIndividuals []dto.Individual
	var copyRightsCorporates []dto.Corporate
	for _, val := range req.Authors.Individuals {
		authorsIndividuals = append(authorsIndividuals, dto.Individual{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthName:    val.AuthName,
		})
	}
	for _, val := range req.Authors.Corporates {
		authorsCorporates = append(authorsCorporates, dto.Corporate{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthName:    val.AuthName,
		})
	}
	for _, val := range req.Copyrights.Individuals {
		copyRightsIndividuals = append(copyRightsIndividuals, dto.Individual{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthName:    val.AuthName,
		})
	}
	for _, val := range req.Copyrights.Corporates {
		copyRightsCorporates = append(copyRightsCorporates, dto.Corporate{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthName:    val.AuthName,
		})
	}
	params := dto.RegisterRequest{
		ProjectID:    authData.ProjectId,
		RegisterType: req.RegisterType,
		OperationID:  req.OperationID,
		UserID:       req.UserID,
		ProductInfo: dto.ProductInfo{
			Name:          req.ProductInfo.Name,
			CatName:       req.ProductInfo.CatName,
			CoverImg:      req.ProductInfo.CoverImg,
			File:          req.ProductInfo.File,
			Description:   req.ProductInfo.Description,
			CreateNatName: req.ProductInfo.CreateNatName,
			CreateTime:    req.ProductInfo.CreateTime,
			CreateAddr:    req.ProductInfo.CreateAddr,
			IsPublished:   req.ProductInfo.IsPublished,
			PubAddr:       req.ProductInfo.PubAddr,
			PubTime:       req.ProductInfo.PubTime,
			PubChannel:    req.ProductInfo.PubChannel,
			PubAnnex:      req.ProductInfo.PubAnnex,
		},
		RightsInfo: dto.RightsInfo{
			Hold:          req.RightsInfo.Hold,
			HoldName:      req.RightsInfo.HoldName,
			HoldExp:       req.RightsInfo.HoldExp,
			RightDocument: req.RightsInfo.RightDocument,
		},
		Authors: dto.Authors{
			Individuals: authorsIndividuals,
			Corporates:  authorsCorporates,
		},
		Copyrights: dto.Copyrights{
			Individuals: copyRightsIndividuals,
			Corporates:  copyRightsCorporates,
		},
		ContactNum:  req.ContactNum,
		Email:       req.Email,
		UrgentTime:  req.UrgentTime,
		CallbackURL: req.CallbackURL,
		AuthFile:    req.AuthFile,
		Metadata:    nil,
	}

	return r.svc.Register(&params)
}

func (r Rights) EditRegister(ctx context.Context, request interface{}) (response interface{}, err error) {
	req := request.(*vo.EditRegisterRequest)
	// 校验参数
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}

	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}

	authData := r.AuthData(ctx)
	var authorsIndividuals []dto.Individual
	var authorsCorporates []dto.Corporate
	var copyRightsIndividuals []dto.Individual
	var copyRightsCorporates []dto.Corporate
	for _, val := range req.Authors.Individuals {
		authorsIndividuals = append(authorsIndividuals, dto.Individual{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthName:    val.AuthName,
		})
	}
	for _, val := range req.Authors.Corporates {
		authorsCorporates = append(authorsCorporates, dto.Corporate{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthName:    val.AuthName,
		})
	}
	for _, val := range req.Copyrights.Individuals {
		copyRightsIndividuals = append(copyRightsIndividuals, dto.Individual{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthName:    val.AuthName,
		})
	}
	for _, val := range req.Copyrights.Corporates {
		copyRightsCorporates = append(copyRightsCorporates, dto.Corporate{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthName:    val.AuthName,
		})
	}
	params := dto.EditRegisterRequest{
		ProjectID:    authData.ProjectId,
		RegisterType: req.RegisterType,
		OperationID:  req.OperationID,
		UserID:       req.UserID,
		ProductInfo: dto.ProductInfo{
			Name:          req.ProductInfo.Name,
			CatName:       req.ProductInfo.CatName,
			CoverImg:      req.ProductInfo.CoverImg,
			File:          req.ProductInfo.File,
			Description:   req.ProductInfo.Description,
			CreateNatName: req.ProductInfo.CreateNatName,
			CreateTime:    req.ProductInfo.CreateTime,
			CreateAddr:    req.ProductInfo.CreateAddr,
			IsPublished:   req.ProductInfo.IsPublished,
			PubAddr:       req.ProductInfo.PubAddr,
			PubTime:       req.ProductInfo.PubTime,
			PubChannel:    req.ProductInfo.PubChannel,
			PubAnnex:      req.ProductInfo.PubAnnex,
		},
		RightsInfo: dto.RightsInfo{
			Hold:          req.RightsInfo.Hold,
			HoldName:      req.RightsInfo.HoldName,
			HoldExp:       req.RightsInfo.HoldExp,
			RightDocument: req.RightsInfo.RightDocument,
		},
		Authors: dto.Authors{
			Individuals: authorsIndividuals,
			Corporates:  authorsCorporates,
		},
		Copyrights: dto.Copyrights{
			Individuals: copyRightsIndividuals,
			Corporates:  copyRightsCorporates,
		},
		ContactNum:  req.ContactNum,
		Email:       req.Email,
		UrgentTime:  req.UrgentTime,
		CallbackURL: req.CallbackURL,
		AuthFile:    req.AuthFile,
		Metadata:    nil,
	}

	return r.svc.EditRegister(&params)
}

func (r Rights) QueryRegister(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.QueryRegisterRequest{
		ProjectID:    authData.ProjectId,
		RegisterType: r.RegisterType(ctx),
		OperationID:  r.OperationID(ctx),
	}

	return r.svc.QueryRegister(&param)
}

func (r Rights) UserAuth(ctx context.Context, request interface{}) (response interface{}, err error) {
	panic("implement me")
}

func (r Rights) EditUserAuth(ctx context.Context, request interface{}) (response interface{}, err error) {
	panic("implement me")
}

func (r Rights) QueryUserAuth(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.QueryUserAuthRequest{
		ProjectID:    authData.ProjectId,
		RegisterType: r.RegisterType(ctx),
		AuthType:     r.AuthType(ctx),
		AuthNum:      r.AuthNum(ctx),
	}

	return r.svc.QueryUserAuth(&param)
}

func (r Rights) Dict(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.DictRequest{
		ProjectID:    authData.ProjectId,
		RegisterType: r.RegisterType(ctx),
	}
	return r.svc.Dict(&param)
}

func (r Rights) Region(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.RegionRequest{
		ProjectID:    authData.ProjectId,
		ParentID:     r.ParentID(ctx),
		RegisterType: r.RegisterType(ctx),
	}
	return r.svc.Region(&param)
}

func (Rights) RegisterType(ctx context.Context) uint64 {
	registerType := ctx.Value("register_type")

	if registerType == 0 {
		return 1
	}
	return registerType.(uint64)
}

func (Rights) ParentID(ctx context.Context) string {
	parentID := ctx.Value("parent_id")

	if parentID == nil {
		return ""
	}
	return parentID.(string)
}

func (Rights) OperationID(ctx context.Context) string {
	operationID := ctx.Value("operation_id")

	if operationID == nil {
		return ""
	}
	return operationID.(string)
}

func (Rights) AuthType(ctx context.Context) string {
	authType := ctx.Value("auth_type")

	if authType == nil {
		return ""
	}
	return authType.(string)
}

func (Rights) AuthNum(ctx context.Context) string {
	authNum := ctx.Value("auth_num")

	if authNum == nil {
		return ""
	}
	return authNum.(string)
}
