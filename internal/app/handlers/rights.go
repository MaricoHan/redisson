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

	Delivery(ctx context.Context, request interface{}) (response interface{}, err error)
	EditDelivery(ctx context.Context, request interface{}) (response interface{}, err error)
	DeliveryInfo(ctx context.Context, request interface{}) (response interface{}, err error)

	Change(ctx context.Context, request interface{}) (response interface{}, err error)
	EditChange(ctx context.Context, request interface{}) (response interface{}, err error)
	ChangeInfo(ctx context.Context, request interface{}) (response interface{}, err error)

	Transfer(ctx context.Context, request interface{}) (response interface{}, err error)
	EditTransfer(ctx context.Context, request interface{}) (response interface{}, err error)
	TransferInfo(ctx context.Context, request interface{}) (response interface{}, err error)

	Revoke(ctx context.Context, request interface{}) (response interface{}, err error)
	EditRevoke(ctx context.Context, request interface{}) (response interface{}, err error)
	RevokeInfo(ctx context.Context, request interface{}) (response interface{}, err error)
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
			Hash:          req.ProductInfo.Hash,
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

	return r.svc.Register(ctx, &params)
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
			Hash:          req.ProductInfo.Hash,
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

	return r.svc.EditRegister(ctx, &params)
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

	return r.svc.QueryRegister(ctx, &param)
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
			ContactNum:      req.AuthInfoCorporate.ContactNum,
			Email:           req.AuthInfoCorporate.Email,
			IndustryCode:    req.AuthInfoCorporate.IndustryCode,
			IndustryName:    req.AuthInfoCorporate.IndustryName,
		},
		CallbackUrl: req.CallbackUrl,
	}

	return r.svc.UserAuth(ctx, &params)
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
			ContactNum:      req.AuthInfoCorporate.ContactNum,
			Email:           req.AuthInfoCorporate.Email,
			IndustryCode:    req.AuthInfoCorporate.IndustryCode,
			IndustryName:    req.AuthInfoCorporate.IndustryName,
		},
		CallbackUrl: req.CallbackUrl,
	}

	return r.svc.EditUserAuth(ctx, &params)
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

	return r.svc.QueryUserAuth(ctx, &param)
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
	return r.svc.Dict(ctx, &param)
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
	return r.svc.Region(ctx, &param)
}

func (r Rights) Delivery(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.DeliveryRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.DeliveryRequest{
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
	return r.svc.Delivery(ctx, &param)
}

func (r Rights) EditDelivery(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.EditDeliveryRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.EditDeliveryRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: req.RegisterType,
		OperationID:  r.OperationID(ctx),
		Addr:         req.Addr,
		Postcode:     req.Postcode,
		Recipient:    req.Recipient,
		PhoneNum:     req.PhoneNum,
	}
	return r.svc.EditDelivery(ctx, &param)
}

func (r Rights) DeliveryInfo(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.DeliveryInfoRequest{
		Code:           authData.Code,
		Module:         authData.Module,
		ProjectID:      authData.ProjectId,
		RegisterType:   r.RegisterType(ctx),
		ProductID:      r.ProductID(ctx),
		CertificateNum: r.CertificateNum(ctx),
	}
	return r.svc.DeliveryInfo(ctx, &param)
}

func (r Rights) Change(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.ChangeRequest)
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

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	copyrighterCorporate := dto.CopyrighterCorporate{
		CompanyName: req.CopyrighterCorporate.CompanyName,
		BusLicImg:   req.CopyrighterCorporate.BusLicImg,
	}

	copyrighterIndividual := dto.CopyrighterIndividual{
		RealName:   req.CopyrighterIndividual.RealName,
		IDCardFimg: req.CopyrighterIndividual.IDCardFimg,
		IDCardBimg: req.CopyrighterIndividual.IDCardBimg,
		IDCardHimg: req.CopyrighterIndividual.IDCardHimg,
	}
	param := dto.ChangeRequest{
		Code:                  authData.Code,
		Module:                authData.Module,
		ProjectID:             authData.ProjectId,
		RegisterType:          req.RegisterType,
		OperationID:           operationId,
		ProductID:             req.ProductID,
		CertificateNum:        req.CertificateNum,
		Name:                  req.Name,
		CatName:               req.CatName,
		CopyrighterCorporate:  copyrighterCorporate,
		CopyrighterIndividual: copyrighterIndividual,
		ProofFiles:            req.ProofFiles,
		UrgentTime:            req.UrgentTime,
	}
	return r.svc.Change(ctx, &param)
}

func (r Rights) EditChange(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.EditChangeRequest)
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

	// 获取账户基本信息
	authData := r.AuthData(ctx)

	copyrighterCorporate := dto.CopyrighterCorporate{
		CompanyName: req.CopyrighterCorporate.CompanyName,
		BusLicImg:   req.CopyrighterCorporate.BusLicImg,
	}

	copyrighterIndividual := dto.CopyrighterIndividual{
		RealName:   req.CopyrighterIndividual.RealName,
		IDCardFimg: req.CopyrighterIndividual.IDCardFimg,
		IDCardBimg: req.CopyrighterIndividual.IDCardBimg,
		IDCardHimg: req.CopyrighterIndividual.IDCardHimg,
	}
	param := dto.EditChangeRequest{
		Code:                  authData.Code,
		Module:                authData.Module,
		ProjectID:             authData.ProjectId,
		RegisterType:          req.RegisterType,
		OperationID:           operationId,
		Name:                  req.Name,
		CatName:               req.CatName,
		CopyrighterCorporate:  copyrighterCorporate,
		CopyrighterIndividual: copyrighterIndividual,
		ProofFiles:            req.ProofFiles,
	}
	return r.svc.EditChange(ctx, &param)
}

func (r Rights) ChangeInfo(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 校验参数
	operationId := strings.TrimSpace(r.OperationID(ctx))
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if r.RegisterType(ctx) == 0 {
		return nil, errors2.New(errors2.ClientParams, "register_type can not be nil")
	}

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.ChangeInfoRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: r.RegisterType(ctx),
		OperationID:  operationId,
	}
	return r.svc.ChangeInfo(ctx, &param)
}

func (r Rights) Transfer(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.TransferRequest)
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

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.TransferRequest{
		Code:             authData.Code,
		Module:           authData.Module,
		ProjectID:        authData.ProjectId,
		RegisterType:     req.RegisterType,
		OperationID:      operationId,
		CertificateNum:   req.CertificateNum,
		ProductID:        req.ProductID,
		AuthorityName:    req.AuthorityName,
		AuthorityIDType:  req.AuthorityIDType,
		AuthorityIDNum:   req.AuthorityIDNum,
		AuthoritedName:   req.AuthoritedIDName,
		AuthoritedIDType: req.AuthoritedIDType,
		AuthoritedIDNum:  req.AuthoritedIDNum,
		AuthInstructions: req.AuthInstructions,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Scope:            req.Scope,
		ContractAmount:   req.ContractAmount,
		ContractFiles:    req.ContractFiles,
		UrgentTime:       req.UrgentTime,
	}
	return r.svc.Transfer(ctx, &param)
}

func (r Rights) EditTransfer(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.EditTransferRequest)
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

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.EditTransferRequest{
		Code:             authData.Code,
		Module:           authData.Module,
		ProjectID:        authData.ProjectId,
		RegisterType:     req.RegisterType,
		OperationID:      operationId,
		AuthorityName:    req.AuthorityName,
		AuthorityIDType:  req.AuthorityIDType,
		AuthorityIDNum:   req.AuthorityIDNum,
		AuthoritedName:   req.AuthoritedIDName,
		AuthoritedIDType: req.AuthoritedIDType,
		AuthoritedIDNum:  req.AuthoritedIDNum,
		AuthInstructions: req.AuthInstructions,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Scope:            req.Scope,
		ContractAmount:   req.ContractAmount,
		ContractFiles:    req.ContractFiles,
	}
	return r.svc.EditTransfer(ctx, &param)
}

func (r Rights) TransferInfo(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 校验参数
	operationId := strings.TrimSpace(r.OperationID(ctx))
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if r.RegisterType(ctx) == 0 {
		return nil, errors2.New(errors2.ClientParams, "register_type can not be nil")
	}

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.TransferInfoRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: r.RegisterType(ctx),
		OperationID:  operationId,
	}
	return r.svc.TransferInfo(ctx, &param)
}

func (r Rights) Revoke(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.RevokeRequest)
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

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.RevokeRequest{
		Code:           authData.Code,
		Module:         authData.Module,
		ProjectID:      authData.ProjectId,
		RegisterType:   req.RegisterType,
		OperationID:    operationId,
		ProductID:      req.ProductID,
		CertificateNum: req.CertificateNum,
	}
	return r.svc.Revoke(ctx, &param)
}

func (r Rights) EditRevoke(ctx context.Context, request interface{}) (response interface{}, err error) {
	req, ok := request.(*vo.EditTransferRequest)
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

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.EditRevokeRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: req.RegisterType,
		OperationID:  operationId,
	}
	return r.svc.EditRevoke(ctx, &param)
}

func (r Rights) RevokeInfo(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 校验参数
	operationId := strings.TrimSpace(r.OperationID(ctx))
	if operationId == "" {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationID)
	}
	if len([]rune(operationId)) == 0 || len([]rune(operationId)) >= 65 {
		return nil, errors2.New(errors2.ClientParams, errors2.ErrOperationIDLen)
	}
	if r.RegisterType(ctx) == 0 {
		return nil, errors2.New(errors2.ClientParams, "register_type can not be nil")
	}

	// 获取账户基本信息
	authData := r.AuthData(ctx)
	param := dto.RevokeInfoRequest{
		Code:         authData.Code,
		Module:       authData.Module,
		ProjectID:    authData.ProjectId,
		RegisterType: r.RegisterType(ctx),
		OperationID:  operationId,
	}
	return r.svc.RevokeInfo(ctx, &param)
}

func (Rights) RegisterType(ctx context.Context) uint64 {
	registerType := ctx.Value("register_type")

	if registerType == nil {
		return 0
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
