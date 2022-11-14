package handlers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	service "gitlab.bianjie.ai/avata/open-api/internal/app/services"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"strconv"
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

	PostCert(ctx context.Context, request interface{}) (response interface{}, err error)
	EditPostCert(ctx context.Context, request interface{}) (response interface{}, err error)
	PostCertInfo(ctx context.Context, request interface{}) (response interface{}, err error)
}

type Rights struct {
	base
	svc service.IRights
}

func NewRights(svc service.IRights) *Rights {
	return &Rights{svc: svc}
}

func (r Rights) Register(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.RegisterRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}

	// 校验参数
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if req.RegisterType == 0 {
		return nil, errors2.New(errors2.ClientParams, "register_type can not be nil")
	}
	if req.UserID == "" {
		return nil, errors2.New(errors2.ClientParams, "user_id can not be nil")
	}
	if req.ContactNum == "" {
		return nil, errors2.New(errors2.ClientParams, "contact_num can not be nil")
	}
	if req.CallbackURL == "" {
		return nil, errors2.New(errors2.ClientParams, "callback_url can not be nil")
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
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range req.Authors.Corporates {
		authorsCorporates = append(authorsCorporates, dto.Corporate{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range req.Copyrights.Individuals {
		copyRightsIndividuals = append(copyRightsIndividuals, dto.Individual{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range req.Copyrights.Corporates {
		copyRightsCorporates = append(copyRightsCorporates, dto.Corporate{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthNum,
		})
	}
	params := dto.RegisterRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: req.RegisterType,
		OperationID:  operationId,
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
	req, ok := request.(*vo.EditRegisterRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 校验参数
	operationId := strings.TrimSpace(r.OperationID(ctx))
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if req.RegisterType == 0 {
		return nil, errors2.New(errors2.ClientParams, "register_type can not be nil")
	}
	if req.UserID == "" {
		return nil, errors2.New(errors2.ClientParams, "user_id can not be nil")
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
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range req.Authors.Corporates {
		authorsCorporates = append(authorsCorporates, dto.Corporate{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range req.Copyrights.Individuals {
		copyRightsIndividuals = append(copyRightsIndividuals, dto.Individual{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range req.Copyrights.Corporates {
		copyRightsCorporates = append(copyRightsCorporates, dto.Corporate{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthNum,
		})
	}
	params := dto.EditRegisterRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: req.RegisterType,
		OperationID:  operationId,
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
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: r.RegisterType(ctx),
		OperationID:  r.OperationID(ctx),
	}

	return r.svc.QueryRegister(&param)
}

func (r Rights) UserAuth(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.UserAuthRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 校验参数
	operationId := strings.TrimSpace(req.OperationID)
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if req.RegisterType == 0 {
		return nil, errors2.New(errors2.ClientParams, "register_type can not be nil")
	}

	authData := r.AuthData(ctx)
	params := dto.UserAuthRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: req.RegisterType,
		OperationID:  operationId,
		AuthType:     req.AuthType,
		AuthInfoIndividual: dto.AuthInfoIndividual{
			RealName:        req.AuthInfoIndividual.RealName,
			IDCardNum:       req.AuthInfoIndividual.IDCardNum,
			IDCardFimg:      req.AuthInfoIndividual.IDCardFimg,
			IDCardBimg:      req.AuthInfoIndividual.IDCardBimg,
			IDCardHimg:      req.AuthInfoIndividual.IDCardHimg,
			IDCardStartDate: req.AuthInfoIndividual.IDCardStartDate,
			IDCardEndDate:   req.AuthInfoIndividual.IDCardEndDate,
			IDCardProvince:  req.AuthInfoIndividual.IDCardProvince,
			IDCardCity:      req.AuthInfoIndividual.IDCardCity,
			IDCardArea:      req.AuthInfoIndividual.IDCardArea,
			ContactNum:      req.AuthInfoIndividual.ContactNum,
			ContactAddr:     req.AuthInfoIndividual.ContactAddr,
			Postcode:        req.AuthInfoIndividual.Postcode,
			Contact:         req.AuthInfoIndividual.Contact,
			Email:           req.AuthInfoIndividual.Email,
			IndustryCode:    req.AuthInfoIndividual.IndustryCode,
			IndustryName:    req.AuthInfoIndividual.IndustryName,
		},
		AuthInfoCorporate: dto.AuthInfoCorporate{
			CardType:        req.AuthInfoCorporate.CardType,
			CompanyName:     req.AuthInfoCorporate.CompanyName,
			BusLicNum:       req.AuthInfoCorporate.BusLicNum,
			CompanyAddr:     req.AuthInfoCorporate.CompanyAddr,
			BusLicImg:       req.AuthInfoCorporate.BusLicImg,
			BusLicStartDate: req.AuthInfoCorporate.BusLicStartDate,
			BusLicEndDate:   req.AuthInfoCorporate.BusLicEndDate,
			BusLicProvince:  req.AuthInfoCorporate.BusLicProvince,
			BusLicCity:      req.AuthInfoCorporate.BusLicCity,
			BusLicArea:      req.AuthInfoCorporate.BusLicArea,
			Postcode:        req.AuthInfoCorporate.Postcode,
			Contact:         req.AuthInfoCorporate.Contact,
			Email:           req.AuthInfoCorporate.Email,
			IndustryCode:    req.AuthInfoCorporate.IndustryCode,
			IndustryName:    req.AuthInfoCorporate.IndustryName,
		},
		CallbackUrl: req.CallbackUrl,
	}

	return r.svc.UserAuth(&params)
}

func (r Rights) EditUserAuth(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.EditUserAuthRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 校验参数
	operationId := strings.TrimSpace(r.OperationID(ctx))
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if req.RegisterType == 0 {
		return nil, errors2.New(errors2.ClientParams, "register_type can not be nil")
	}

	authData := r.AuthData(ctx)
	params := dto.EditUserAuthRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: req.RegisterType,
		OperationID:  operationId,
		AuthType:     req.AuthType,
		AuthInfoIndividual: dto.AuthInfoIndividual{
			RealName:        req.AuthInfoIndividual.RealName,
			IDCardNum:       req.AuthInfoIndividual.IDCardNum,
			IDCardFimg:      req.AuthInfoIndividual.IDCardFimg,
			IDCardBimg:      req.AuthInfoIndividual.IDCardBimg,
			IDCardHimg:      req.AuthInfoIndividual.IDCardHimg,
			IDCardStartDate: req.AuthInfoIndividual.IDCardStartDate,
			IDCardEndDate:   req.AuthInfoIndividual.IDCardEndDate,
			IDCardProvince:  req.AuthInfoIndividual.IDCardProvince,
			IDCardCity:      req.AuthInfoIndividual.IDCardCity,
			IDCardArea:      req.AuthInfoIndividual.IDCardArea,
			ContactNum:      req.AuthInfoIndividual.Contact,
			ContactAddr:     req.AuthInfoIndividual.ContactAddr,
			Postcode:        req.AuthInfoIndividual.Postcode,
			Contact:         req.AuthInfoIndividual.ContactAddr,
			Email:           req.AuthInfoIndividual.Email,
			IndustryCode:    req.AuthInfoIndividual.IndustryCode,
			IndustryName:    req.AuthInfoIndividual.IndustryName,
		},
		AuthInfoCorporate: dto.AuthInfoCorporate{
			CardType:        req.AuthInfoCorporate.CardType,
			CompanyName:     req.AuthInfoCorporate.CompanyName,
			BusLicNum:       req.AuthInfoCorporate.BusLicNum,
			CompanyAddr:     req.AuthInfoCorporate.CompanyAddr,
			BusLicImg:       req.AuthInfoCorporate.BusLicImg,
			BusLicStartDate: req.AuthInfoCorporate.BusLicStartDate,
			BusLicEndDate:   req.AuthInfoCorporate.BusLicEndDate,
			BusLicProvince:  req.AuthInfoCorporate.BusLicProvince,
			BusLicCity:      req.AuthInfoCorporate.BusLicCity,
			BusLicArea:      req.AuthInfoCorporate.BusLicArea,
			Postcode:        req.AuthInfoCorporate.Postcode,
			Contact:         req.AuthInfoCorporate.Contact,
			Email:           req.AuthInfoCorporate.Email,
			IndustryCode:    req.AuthInfoCorporate.IndustryCode,
			IndustryName:    req.AuthInfoCorporate.IndustryName,
		},
		CallbackUrl: req.CallbackUrl,
	}

	return r.svc.EditUserAuth(&params)
}

func (r Rights) QueryUserAuth(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.QueryUserAuthRequest{
		Code:         authData.Code,
		Module:       authData.Module,
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
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: r.RegisterType(ctx),
	}
	return r.svc.Dict(&param)
}

func (r Rights) Region(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.RegionRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		ParentID:     r.ParentID(ctx),
		RegisterType: r.RegisterType(ctx),
	}
	return r.svc.Region(&param)
}

func (r Rights) PostCert(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.PostCertRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.PostCertRequest{
		Code:           authData.Code,
		Module:         authData.Module,
		ProjectID:      authData.ProjectId,
		RegisterType:   req.RegisterType,
		OperationID:    req.OperationID,
		ProductID:      req.ProductID,
		CertificateNum: req.CertificateNum,
		Addr:           req.Addr,
		Postcode:       req.Postcode,
		Recipient:      req.Recipient,
		PhoneNum:       req.PhoneNum,
	}
	return r.svc.PostCert(&param)
}

func (r Rights) EditPostCert(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.EditPostCertRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.EditPostCertRequest{
		Code:           authData.Code,
		Module:         authData.Module,
		ProjectID:      authData.ProjectId,
		RegisterType:   req.RegisterType,
		OperationID:    r.OperationID(ctx),
		ProductID:      req.ProductID,
		CertificateNum: req.CertificateNum,
		Addr:           req.Addr,
		Postcode:       req.Postcode,
		Recipient:      req.Recipient,
		PhoneNum:       req.PhoneNum,
	}
	return r.svc.EditPostCert(&param)
}

func (r Rights) PostCertInfo(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.PostCertInfoRequest{
		Code:           authData.Code,
		Module:         authData.Module,
		ProjectID:      authData.ProjectId,
		RegisterType:   r.RegisterType(ctx),
		ProductID:      r.ProductID(ctx),
		CertificateNum: r.CertificateNum(ctx),
	}
	return r.svc.PostCertInfo(&param)
}

func (Rights) RegisterType(ctx context.Context) uint64 {
	registerType := ctx.Value("register_type")

	if registerType == 0 {
		return 1
	}
	r := registerType.(string)
	parseUint, _ := strconv.ParseUint(r, 10, 64)
	return parseUint
}

func (Rights) ParentID(ctx context.Context) uint64 {
	parentID := ctx.Value("parent_id")

	if parentID == nil {
		return 0
	}
	p := parentID.(string)
	parseUint, _ := strconv.ParseUint(p, 10, 64)
	return parseUint
}

func (Rights) OperationID(ctx context.Context) string {
	operationID := ctx.Value("operation_id")

	if operationID == nil {
		return ""
	}
	return operationID.(string)
}

func (Rights) AuthType(ctx context.Context) uint32 {
	authType := ctx.Value("auth_type")

	if authType == nil {
		return 0
	}
	p := authType.(string)
	parseUint, _ := strconv.ParseUint(p, 10, 64)
	return uint32(parseUint)
}

func (Rights) AuthNum(ctx context.Context) string {
	authNum := ctx.Value("auth_num")

	if authNum == nil {
		return ""
	}
	return authNum.(string)
}

func (Rights) CertificateNum(ctx context.Context) string {
	certificateNum := ctx.Value("certificate_num")

	if certificateNum == nil {
		return ""
	}
	return certificateNum.(string)
}

func (Rights) ProductID(ctx context.Context) string {
	productID := ctx.Value("product_id")

	if productID == nil {
		return ""
	}
	return productID.(string)
}
