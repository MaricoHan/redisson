package services

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"gitlab.bianjie.ai/avata/services/api/pb/rights"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"time"
)

type IRights interface {
	Register(params *dto.RegisterRequest) (*dto.RegisterResponse, error)
	EditRegister(params *dto.EditRegisterRequest) (*dto.EditRegisterResponse, error)
	QueryRegister(params *dto.QueryRegisterRequest) (*dto.QueryRegisterResponse, error)
	UserAuth(params *dto.UserAuthRequest) (*dto.UserAuthResponse, error)
	EditUserAuth(params *dto.EditUserAuthRequest) (*dto.EditUserAuthResponse, error)
	QueryUserAuth(params *dto.QueryUserAuthRequest) (*dto.QueryUserAuthResponse, error)

	Dict(params *dto.DictRequest) (*dto.DictResponse, error)
	Region(params *dto.RegionRequest) (*dto.RegionResponse, error)
}

type Rights struct {
	logger *log.Entry
}

func NewRights(logger *log.Logger) *Rights {
	return &Rights{
		logger: logger.WithField("service", "rights"),
	}
}

func (r Rights) Register(params *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Register")

	var authorsIndividual []*rights.Person
	var authorsCorporate []*rights.Company
	var copyrightsIndividual []*rights.Person
	var copyrightsCorporate []*rights.Company
	for _, val := range params.Authors.Individuals {
		authorsIndividual = append(authorsIndividual, &rights.Person{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthName,
		})
	}
	for _, val := range params.Authors.Corporates {
		authorsCorporate = append(authorsCorporate, &rights.Company{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthName,
		})
	}
	for _, val := range params.Copyrights.Individuals {
		copyrightsIndividual = append(copyrightsIndividual, &rights.Person{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthName,
		})
	}
	for _, val := range params.Copyrights.Corporates {
		copyrightsCorporate = append(copyrightsCorporate, &rights.Company{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthName,
		})
	}

	req := rights.RegisterRequest{
		Code:        "",
		Module:      "",
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
		ProductInfo: &rights.ProductInfo{
			Name:          params.ProductInfo.Name,
			CatName:       params.ProductInfo.CatName,
			CoverImg:      params.ProductInfo.CoverImg,
			File:          params.ProductInfo.File,
			Description:   params.ProductInfo.Description,
			CreateNatName: params.ProductInfo.CreateNatName,
			CreateTime:    params.ProductInfo.CreateTime,
			CreateAddr:    params.ProductInfo.CreateAddr,
			IsPublished:   params.ProductInfo.IsPublished,
			PubAddr:       params.ProductInfo.PubAddr,
			PubTime:       params.ProductInfo.PubTime,
			PubChannel:    params.ProductInfo.PubChannel,
			PubAnnex:      params.ProductInfo.PubAnnex,
		},
		RightsInfo: &rights.RightsInfo{
			Hold:          params.RightsInfo.Hold,
			HoldName:      params.RightsInfo.HoldName,
			HoldExp:       params.RightsInfo.HoldExp,
			RightDocument: params.RightsInfo.RightDocument,
		},
		Authors: &rights.Authors{
			AuthorsIndividual: authorsIndividual,
			AuthorsCorporate:  authorsCorporate,
		},
		Copyrights: &rights.Copyrights{
			CopyrightsIndividual: copyrightsIndividual,
			CopyrightsCorporate:  copyrightsCorporate,
		},
		ContactNum:  params.ContactNum,
		Email:       params.Email,
		UrgentTime:  params.UrgentTime,
		CallbackUrl: params.CallbackURL,
		AuthFile:    params.AuthFile,
		Metadata:    "",
	}

	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors.New("") // todo
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.Register(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, errors.New("") // todo
	}

	return &dto.RegisterResponse{OperationID: resp.OperationId}, nil
}

func (r Rights) EditRegister(params *dto.EditRegisterRequest) (*dto.EditRegisterResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "EditRegister")

	var authorsIndividual []*rights.Person
	var authorsCorporate []*rights.Company
	var copyrightsIndividual []*rights.Person
	var copyrightsCorporate []*rights.Company
	for _, val := range params.Authors.Individuals {
		authorsIndividual = append(authorsIndividual, &rights.Person{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthName,
		})
	}
	for _, val := range params.Authors.Corporates {
		authorsCorporate = append(authorsCorporate, &rights.Company{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthName,
		})
	}
	for _, val := range params.Copyrights.Individuals {
		copyrightsIndividual = append(copyrightsIndividual, &rights.Person{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthName,
		})
	}
	for _, val := range params.Copyrights.Corporates {
		copyrightsCorporate = append(copyrightsCorporate, &rights.Company{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthName,
		})
	}

	req := rights.RegisterRequest{
		Code:        "",
		Module:      "",
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
		ProductInfo: &rights.ProductInfo{
			Name:          params.ProductInfo.Name,
			CatName:       params.ProductInfo.CatName,
			CoverImg:      params.ProductInfo.CoverImg,
			File:          params.ProductInfo.File,
			Description:   params.ProductInfo.Description,
			CreateNatName: params.ProductInfo.CreateNatName,
			CreateTime:    params.ProductInfo.CreateTime,
			CreateAddr:    params.ProductInfo.CreateAddr,
			IsPublished:   params.ProductInfo.IsPublished,
			PubAddr:       params.ProductInfo.PubAddr,
			PubTime:       params.ProductInfo.PubTime,
			PubChannel:    params.ProductInfo.PubChannel,
			PubAnnex:      params.ProductInfo.PubAnnex,
		},
		RightsInfo: &rights.RightsInfo{
			Hold:          params.RightsInfo.Hold,
			HoldName:      params.RightsInfo.HoldName,
			HoldExp:       params.RightsInfo.HoldExp,
			RightDocument: params.RightsInfo.RightDocument,
		},
		Authors: &rights.Authors{
			AuthorsIndividual: authorsIndividual,
			AuthorsCorporate:  authorsCorporate,
		},
		Copyrights: &rights.Copyrights{
			CopyrightsIndividual: copyrightsIndividual,
			CopyrightsCorporate:  copyrightsCorporate,
		},
		ContactNum:  params.ContactNum,
		Email:       params.Email,
		UrgentTime:  params.UrgentTime,
		CallbackUrl: params.CallbackURL,
		AuthFile:    params.AuthFile,
		Metadata:    "",
	}

	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors.New("") // todo
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.EditRegister(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, errors.New("") // todo
	}

	return &dto.EditRegisterResponse{OperationID: resp.OperationId}, nil
}

func (r Rights) QueryRegister(params *dto.QueryRegisterRequest) (*dto.QueryRegisterResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "QueryRegister")

	req := rights.RegisterInfoRequest{OperationId: params.OperationID}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors.New("") // todo
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.RegisterInfo(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, errors.New("") // todo
	}
	result := &dto.QueryRegisterResponse{
		OperationID:       resp.OperationId,
		AuditStatus:       resp.AuditStatus,
		AuditFile:         resp.AuditFile,
		AuditOpinion:      resp.AuditOpinion,
		CertificateStatus: resp.CertificateStatus,
		CertificateNum:    resp.CertificateNum,
		CertificateURL:    resp.CertificateUrl,
		BackTag:           resp.BackTag,
	}
	return result, nil
}

func (r Rights) UserAuth(params *dto.UserAuthRequest) (*dto.UserAuthResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "UserAuth")

	req := rights.UserAuthRequest{
		Code:        "",
		Module:      "",
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
		AuthType:    params.AuthType,
		AuthInfoIndividual: &rights.AuthPerson{
			RealName:        params.AuthInfoIndividual.RealName,
			IdcardNum:       params.AuthInfoIndividual.IDCardNum,
			IdcardFimg:      params.AuthInfoIndividual.IDCardFimg,
			IdcardBimg:      params.AuthInfoIndividual.IDCardBimg,
			IdcardHimg:      params.AuthInfoIndividual.IDCardHimg,
			IdcardStartDate: params.AuthInfoIndividual.IDCardStartDate,
			IdcardEndDate:   params.AuthInfoIndividual.IDCardEndDate,
			IdcardProvince:  params.AuthInfoIndividual.IDCardProvince,
			IdcardCity:      params.AuthInfoIndividual.IDCardCity,
			IdcardArea:      params.AuthInfoIndividual.IDCardArea,
			ContactNum:      params.AuthInfoIndividual.ContactNum,
			ContactAddr:     params.AuthInfoIndividual.ContactAddr,
			Postcode:        params.AuthInfoIndividual.Postcode,
			Contact:         params.AuthInfoIndividual.Contact,
			Email:           params.AuthInfoIndividual.Email,
			IndustryCode:    params.AuthInfoIndividual.IndustryCode,
			IndustryName:    params.AuthInfoIndividual.IndustryName,
		},
		AuthInfoCorporate: &rights.AuthCompany{
			CardType:        params.AuthInfoCorporate.CardType,
			CompanyName:     params.AuthInfoCorporate.CompanyName,
			BusLicNum:       params.AuthInfoCorporate.BusLicNum,
			CompanyAddr:     params.AuthInfoCorporate.CompanyAddr,
			BusLicImg:       params.AuthInfoCorporate.BusLicImg,
			BusLicStartDate: params.AuthInfoCorporate.BusLicStartDate,
			BusLicEndDate:   params.AuthInfoCorporate.BusLicEndDate,
			BusLicProvince:  params.AuthInfoCorporate.BusLicProvince,
			BusLicCity:      params.AuthInfoCorporate.BusLicCity,
			BusLicArea:      params.AuthInfoCorporate.BusLicArea,
			Postcode:        params.AuthInfoCorporate.Postcode,
			Contact:         params.AuthInfoCorporate.Contact,
			ContactNum:      params.AuthInfoCorporate.Contact,
			Email:           params.AuthInfoCorporate.Email,
			IndustryCode:    params.AuthInfoCorporate.IndustryCode,
			IndustryName:    params.AuthInfoCorporate.IndustryName,
		},
		CallbackUrl: params.CallbackUrl,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors.New("") // todo
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.UserAuth(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, errors.New("") // todo
	}
	result := &dto.UserAuthResponse{
		OperationID:      resp.OperationId,
		UserID:           resp.UserId,
		AuditStatus:      resp.AuditStatus,
		AuditInstruction: resp.AuditInstruction,
	}
	return result, nil
}

func (r Rights) EditUserAuth(params *dto.EditUserAuthRequest) (*dto.EditUserAuthResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "EditUserAuth")

	req := rights.UserAuthRequest{
		Code:        "",
		Module:      "",
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
		AuthType:    params.AuthType,
		AuthInfoIndividual: &rights.AuthPerson{
			RealName:        params.AuthInfoIndividual.RealName,
			IdcardNum:       params.AuthInfoIndividual.IDCardNum,
			IdcardFimg:      params.AuthInfoIndividual.IDCardFimg,
			IdcardBimg:      params.AuthInfoIndividual.IDCardBimg,
			IdcardHimg:      params.AuthInfoIndividual.IDCardHimg,
			IdcardStartDate: params.AuthInfoIndividual.IDCardStartDate,
			IdcardEndDate:   params.AuthInfoIndividual.IDCardEndDate,
			IdcardProvince:  params.AuthInfoIndividual.IDCardProvince,
			IdcardCity:      params.AuthInfoIndividual.IDCardCity,
			IdcardArea:      params.AuthInfoIndividual.IDCardArea,
			ContactNum:      params.AuthInfoIndividual.ContactNum,
			ContactAddr:     params.AuthInfoIndividual.ContactAddr,
			Postcode:        params.AuthInfoIndividual.Postcode,
			Contact:         params.AuthInfoIndividual.Contact,
			Email:           params.AuthInfoIndividual.Email,
			IndustryCode:    params.AuthInfoIndividual.IndustryCode,
			IndustryName:    params.AuthInfoIndividual.IndustryName,
		},
		AuthInfoCorporate: &rights.AuthCompany{
			CardType:        params.AuthInfoCorporate.CardType,
			CompanyName:     params.AuthInfoCorporate.CompanyName,
			BusLicNum:       params.AuthInfoCorporate.BusLicNum,
			CompanyAddr:     params.AuthInfoCorporate.CompanyAddr,
			BusLicImg:       params.AuthInfoCorporate.BusLicImg,
			BusLicStartDate: params.AuthInfoCorporate.BusLicStartDate,
			BusLicEndDate:   params.AuthInfoCorporate.BusLicEndDate,
			BusLicProvince:  params.AuthInfoCorporate.BusLicProvince,
			BusLicCity:      params.AuthInfoCorporate.BusLicCity,
			BusLicArea:      params.AuthInfoCorporate.BusLicArea,
			Postcode:        params.AuthInfoCorporate.Postcode,
			Contact:         params.AuthInfoCorporate.Contact,
			ContactNum:      params.AuthInfoCorporate.Contact,
			Email:           params.AuthInfoCorporate.Email,
			IndustryCode:    params.AuthInfoCorporate.IndustryCode,
			IndustryName:    params.AuthInfoCorporate.IndustryName,
		},
		CallbackUrl: params.CallbackUrl,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors.New("") // todo
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.EditUserAuth(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, errors.New("") // todo
	}
	result := &dto.EditUserAuthResponse{
		Data: resp.Data,
	}
	return result, nil
}

func (r Rights) QueryUserAuth(params *dto.QueryUserAuthRequest) (*dto.QueryUserAuthResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "QueryUserAuth")

	req := rights.UserAuthInfoRequest{
		AuthType: params.AuthType,
		AuthNum:  params.AuthNum,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors.New("") // todo
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.UserAuthInfo(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, errors.New("") // todo
	}
	result := &dto.QueryUserAuthResponse{
		UserID:           resp.UserId,
		AuditStatus:      resp.AuditStatus,
		AuditInstruction: resp.AuditInstruction,
	}
	return result, nil
}

func (r Rights) Dict(params *dto.DictRequest) (*dto.DictResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Dict")
	req := rights.DictRequest{}

	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors.New("") // todo
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.Dict(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, errors.New("") // todo
	}

	result := &dto.DictResponse{
		ProCat:       nil,
		ProCreateNat: nil,
		IndustryCode: nil,
		AutHold:      nil,
	}
	for _, val := range resp.ProCat {
		result.ProCat = append(result.ProCat, dto.KeyValueDetail{
			Key:    val.Key,
			Value:  val.Value,
			Detail: val.Detail,
		})
	}
	for _, val := range resp.ProCreateNat {
		result.ProCreateNat = append(result.ProCreateNat, dto.KeyValueDetail{
			Key:    val.Key,
			Value:  val.Value,
			Detail: val.Detail,
		})
	}
	for _, val := range resp.IndustryCode {
		result.IndustryCode = append(result.IndustryCode, dto.KeyValue{
			Key:   val.Key,
			Value: val.Value,
		})
	}
	for _, val := range resp.AutHold {
		result.AutHold = append(result.AutHold, dto.KeyValue{
			Key:   val.Key,
			Value: val.Value,
		})
	}

	return result, nil
}

func (r Rights) Region(params *dto.RegionRequest) (*dto.RegionResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Region")

	req := rights.RegionRequest{ParentId: params.ParentID}

	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors.New("") // todo
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.Region(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, errors.New("") // todo
	}

	result := &dto.RegionResponse{
		Data: nil,
	}
	for _, val := range resp.Data {
		result.Data = append(result.Data, dto.Region{
			ID:         val.Id,
			Name:       val.Name,
			ParentID:   val.ParentId,
			ShortName:  val.ShortName,
			MergerName: val.MergerName,
			PinYin:     val.Pinyin,
		})
	}

	return result, nil
}
