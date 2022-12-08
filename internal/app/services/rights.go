package services

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"gitlab.bianjie.ai/avata/services/api/pb/rights"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"time"
)

type IRights interface {
	Register(ctx context.Context, params *dto.RegisterRequest) (*dto.RegisterResponse, error)
	EditRegister(ctx context.Context, params *dto.EditRegisterRequest) (*dto.EditRegisterResponse, error)
	QueryRegister(ctx context.Context, params *dto.QueryRegisterRequest) (*dto.QueryRegisterResponse, error)
	UserAuth(ctx context.Context, params *dto.UserAuthRequest) (*dto.UserAuthResponse, error)
	EditUserAuth(ctx context.Context, params *dto.EditUserAuthRequest) (*dto.EditUserAuthResponse, error)
	QueryUserAuth(ctx context.Context, params *dto.QueryUserAuthRequest) (*dto.QueryUserAuthResponse, error)

	Dict(ctx context.Context, params *dto.DictRequest) (*dto.DictResponse, error)
	Region(ctx context.Context, params *dto.RegionRequest) (*dto.RegionResponse, error)

	Delivery(ctx context.Context, params *dto.DeliveryRequest) (*dto.DeliveryResponse, error)
	EditDelivery(ctx context.Context, params *dto.EditDeliveryRequest) (*dto.EditDeliveryResponse, error)
	DeliveryInfo(ctx context.Context, params *dto.DeliveryInfoRequest) (*dto.DeliveryInfoResponse, error)

	Change(ctx context.Context, params *dto.ChangeRequest) (*dto.ChangeResponse, error)
	EditChange(ctx context.Context, params *dto.EditChangeRequest) (*dto.EditChangeResponse, error)
	ChangeInfo(ctx context.Context, params *dto.ChangeInfoRequest) (*dto.ChangeInfoResponse, error)

	Transfer(ctx context.Context, params *dto.TransferRequest) (*dto.TransferResponse, error)
	EditTransfer(ctx context.Context, params *dto.EditTransferRequest) (*dto.EditTransferResponse, error)
	TransferInfo(ctx context.Context, params *dto.TransferInfoRequest) (*dto.TransferInfoResponse, error)

	Revoke(ctx context.Context, params *dto.RevokeRequest) (*dto.RevokeResponse, error)
	EditRevoke(ctx context.Context, params *dto.EditRevokeRequest) (*dto.EditRevokeResponse, error)
	RevokeInfo(ctx context.Context, params *dto.RevokeInfoRequest) (*dto.RevokeInfoResponse, error)
}

type Rights struct {
	logger *log.Entry
}

func NewRights(logger *log.Logger) *Rights {
	return &Rights{
		logger: logger.WithField("service", "rights"),
	}
}

func (r Rights) Register(ctx context.Context, params *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Register")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	var authorsIndividual []*rights.Person
	var authorsCorporate []*rights.Company
	var copyrightsIndividual []*rights.Person
	var copyrightsCorporate []*rights.Company
	for _, val := range params.Authors.Individuals {
		authorsIndividual = append(authorsIndividual, &rights.Person{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range params.Authors.Corporates {
		authorsCorporate = append(authorsCorporate, &rights.Company{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range params.Copyrights.Individuals {
		copyrightsIndividual = append(copyrightsIndividual, &rights.Person{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range params.Copyrights.Corporates {
		copyrightsCorporate = append(copyrightsCorporate, &rights.Company{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthNum,
		})
	}

	req := rights.RegisterRequest{
		Code:        params.Code,
		Module:      params.Module,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
		AuthUserId:  params.UserID,
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
			Hash:          params.ProductInfo.Hash,
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
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	resp, err := grpcClient.Register(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.RegisterResponse{OperationID: resp.OperationId}, nil
}

func (r Rights) EditRegister(ctx context.Context, params *dto.EditRegisterRequest) (*dto.EditRegisterResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "EditRegister")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	var authorsIndividual []*rights.Person
	var authorsCorporate []*rights.Company
	var copyrightsIndividual []*rights.Person
	var copyrightsCorporate []*rights.Company
	for _, val := range params.Authors.Individuals {
		authorsIndividual = append(authorsIndividual, &rights.Person{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range params.Authors.Corporates {
		authorsCorporate = append(authorsCorporate, &rights.Company{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range params.Copyrights.Individuals {
		copyrightsIndividual = append(copyrightsIndividual, &rights.Person{
			IsApplicant: val.IsApplicant,
			RealName:    val.RealName,
			AuthNum:     val.AuthNum,
		})
	}
	for _, val := range params.Copyrights.Corporates {
		copyrightsCorporate = append(copyrightsCorporate, &rights.Company{
			IsApplicant: val.IsApplicant,
			CardType:    val.CardType,
			CompanyName: val.CompanyName,
			AuthNum:     val.AuthNum,
		})
	}

	req := rights.RegisterRequest{
		Code:        params.Code,
		Module:      params.Module,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
		AuthUserId:  params.UserID,
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
			Hash:          params.ProductInfo.Hash,
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
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	resp, err := grpcClient.EditRegister(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.EditRegisterResponse{OperationID: resp.OperationId}, nil
}

func (r Rights) QueryRegister(ctx context.Context, params *dto.QueryRegisterRequest) (*dto.QueryRegisterResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "QueryRegister")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := rights.RegisterInfoRequest{OperationId: params.OperationID, ProjectId: params.ProjectID}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	resp, err := grpcClient.RegisterInfo(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.QueryRegisterResponse{
		OperationID:       resp.OperationId,
		AuditStatus:       resp.AuditStatus,
		AuditFile:         resp.AuditFile,
		AuditOpinion:      resp.AuditOpinion,
		CertificateStatus: resp.CertificateStatus,
		CertificateNum:    resp.CertificateNum,
		ProductID:         resp.ProductId,
		CertificateURL:    resp.CertificateUrl,
		BackTag:           resp.BackTag,
	}
	return result, nil
}

func (r Rights) UserAuth(ctx context.Context, params *dto.UserAuthRequest) (*dto.UserAuthResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "UserAuth")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := rights.UserAuthRequest{
		Code:        params.Code,
		Module:      params.Module,
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
			ContactNum:      params.AuthInfoCorporate.ContactNum,
			Email:           params.AuthInfoCorporate.Email,
			IndustryCode:    params.AuthInfoCorporate.IndustryCode,
			IndustryName:    params.AuthInfoCorporate.IndustryName,
		},
		CallbackUrl: params.CallbackUrl,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	resp, err := grpcClient.UserAuth(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.UserAuthResponse{
		OperationID:      resp.OperationId,
		UserID:           resp.AuthUserId,
		AuditStatus:      resp.AuditStatus,
		AuditInstruction: resp.AuditInstruction,
	}
	return result, nil
}

func (r Rights) EditUserAuth(ctx context.Context, params *dto.EditUserAuthRequest) (*dto.EditUserAuthResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "EditUserAuth")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := rights.UserAuthRequest{
		Code:        params.Code,
		Module:      params.Module,
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
			ContactNum:      params.AuthInfoCorporate.ContactNum,
			Email:           params.AuthInfoCorporate.Email,
			IndustryCode:    params.AuthInfoCorporate.IndustryCode,
			IndustryName:    params.AuthInfoCorporate.IndustryName,
		},
		CallbackUrl: params.CallbackUrl,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	resp, err := grpcClient.EditUserAuth(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.EditUserAuthResponse{
		Data: resp.Data,
	}
	return result, nil
}

func (r Rights) QueryUserAuth(ctx context.Context, params *dto.QueryUserAuthRequest) (*dto.QueryUserAuthResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "QueryUserAuth")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := rights.UserAuthInfoRequest{
		AuthType: params.AuthType,
		AuthNum:  params.AuthNum,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	resp, err := grpcClient.UserAuthInfo(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	result := &dto.QueryUserAuthResponse{
		UserID:           resp.AuthUserId,
		AuditStatus:      resp.AuditStatus,
		AuditInstruction: resp.AuditInstruction,
	}
	return result, nil
}

func (r Rights) Dict(ctx context.Context, params *dto.DictRequest) (*dto.DictResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Dict")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := rights.DictRequest{}

	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	resp, err := grpcClient.Dict(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
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

func (r Rights) Region(ctx context.Context, params *dto.RegionRequest) (*dto.RegionResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Region")
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()

	req := rights.RegionRequest{ParentId: params.ParentID}

	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	resp, err := grpcClient.Region(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
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

func (r Rights) Delivery(ctx context.Context, params *dto.DeliveryRequest) (*dto.DeliveryResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Delivery")

	req := rights.DeliveryRequest{
		Code:           params.Code,
		Module:         params.Module,
		ProjectId:      params.ProjectID,
		OperationId:    params.OperationID,
		ProductId:      params.ProductID,
		CertificateNum: params.CertificateNum,
		Addr:           params.Addr,
		Postcode:       params.Postcode,
		Recipient:      params.Recipient,
		PhoneNum:       params.PhoneNum,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.Delivery(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.DeliveryResponse{OperationID: resp.OperationId}, nil
}

func (r Rights) EditDelivery(ctx context.Context, params *dto.EditDeliveryRequest) (*dto.EditDeliveryResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "EditDelivery")

	req := rights.EditDeliveryRequest{
		Code:        params.Code,
		Module:      params.Module,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
		Addr:        params.Addr,
		Postcode:    params.Postcode,
		Recipient:   params.Recipient,
		PhoneNum:    params.PhoneNum,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.EditDelivery(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.EditDeliveryResponse{}, nil
}

func (r Rights) DeliveryInfo(ctx context.Context, params *dto.DeliveryInfoRequest) (*dto.DeliveryInfoResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "DeliveryInfo")

	req := rights.DeliveryInfoRequest{
		Code:           params.Code,
		Module:         params.Module,
		ProjectId:      params.ProjectID,
		ProductId:      params.ProductID,
		CertificateNum: params.CertificateNum,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.DeliveryInfo(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.DeliveryInfoResponse{
		ProductID:      resp.ProductId,
		CertificateNum: resp.CertificateNum,
		Addr:           resp.Addr,
		Postcode:       resp.Postcode,
		Recipient:      resp.Recipient,
		PhoneNum:       resp.PhoneNum,
		ExpressNum:     resp.ExpressNum,
		Status:         resp.Status,
	}, nil
}

func (r Rights) Change(ctx context.Context, params *dto.ChangeRequest) (*dto.ChangeResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Change")

	req := rights.ChangeRequest{
		Code:            params.Code,
		Module:          params.Module,
		ProjectId:       params.ProjectID,
		OperationId:     params.OperationID,
		ProductId:       params.ProductID,
		CertificateNum:  params.CertificateNum,
		Name:            params.Name,
		CatName:         params.CatName,
		CopyrighterType: params.CopyrighterType,
		CopyrighterName: params.CopyrighterName,
		CardImg1:        params.CardImg1,
		CardImg2:        params.CardImg2,
		CardImg3:        params.CardImg3,
		ProofFiles:      params.ProofFiles,
		UrgentTime:      params.UrgentTime,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.Change(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.ChangeResponse{OperationID: resp.OperationId}, nil
}
func (r Rights) EditChange(ctx context.Context, params *dto.EditChangeRequest) (*dto.EditChangeResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "EditChange")

	req := rights.EditChangeRequest{
		Code:            params.Code,
		Module:          params.Module,
		ProjectId:       params.ProjectID,
		OperationId:     params.OperationID,
		Name:            params.Name,
		CatName:         params.CatName,
		CopyrighterType: params.CopyrighterType,
		CopyrighterName: params.CopyrighterName,
		CardImg1:        params.CardImg1,
		CardImg2:        params.CardImg2,
		CardImg3:        params.CardImg3,
		ProofFiles:      params.ProofFiles,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.EditChange(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.EditChangeResponse{}, nil
}
func (r Rights) ChangeInfo(ctx context.Context, params *dto.ChangeInfoRequest) (*dto.ChangeInfoResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "ChangeInfo")

	req := rights.ChangeInfoRequest{
		Code:        params.Code,
		Module:      params.Module,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.ChangeInfo(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.ChangeInfoResponse{
		ProductId:            resp.ProductId,
		CertificateNum:       resp.CertificateNum,
		Status:               resp.Status,
		ChangeCertificateUrl: resp.ChangeCertificateUrl,
		ErrorMessage:         resp.ErrorMessage,
		ChangeCertificateNum: resp.ChangeCertificateNum,
	}, nil
}

func (r Rights) Transfer(ctx context.Context, params *dto.TransferRequest) (*dto.TransferResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Transfer")

	req := rights.TransferRequest{
		Code:             params.Code,
		Module:           params.Module,
		ProjectId:        params.ProjectID,
		OperationId:      params.OperationID,
		ProductId:        params.ProductID,
		CertificateNum:   params.CertificateNum,
		AuthorityName:    params.AuthorityName,
		AuthorityIdType:  params.AuthorityIDType,
		AuthorityIdNum:   params.AuthorityIDNum,
		AuthoritedName:   params.AuthoritedName,
		AuthoritedIdType: params.AuthoritedIDType,
		AuthoritedIdNum:  params.AuthoritedIDNum,
		AuthInstructions: params.AuthInstructions,
		StartTime:        params.StartTime,
		EndTime:          params.EndTime,
		Scope:            params.Scope,
		ContractAmount:   params.ContractAmount,
		ContractFiles:    params.ContractFiles,
		UrgentTime:       params.UrgentTime,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.Transfer(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.TransferResponse{OperationID: resp.OperationId}, nil
}
func (r Rights) EditTransfer(ctx context.Context, params *dto.EditTransferRequest) (*dto.EditTransferResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "EditTransfer")

	req := rights.EditTransferRequest{
		Code:             params.Code,
		Module:           params.Module,
		ProjectId:        params.ProjectID,
		OperationId:      params.OperationID,
		AuthorityName:    params.AuthorityName,
		AuthorityIdType:  params.AuthorityIDType,
		AuthorityIdNum:   params.AuthorityIDNum,
		AuthoritedName:   params.AuthoritedName,
		AuthoritedIdType: params.AuthoritedIDType,
		AuthoritedIdNum:  params.AuthoritedIDNum,
		AuthInstructions: params.AuthInstructions,
		StartTime:        params.StartTime,
		EndTime:          params.EndTime,
		Scope:            params.Scope,
		ContractAmount:   params.ContractAmount,
		ContractFiles:    params.ContractFiles,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.EditTransfer(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.EditTransferResponse{}, nil
}
func (r Rights) TransferInfo(ctx context.Context, params *dto.TransferInfoRequest) (*dto.TransferInfoResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "TransferInfo")

	req := rights.TransferInfoRequest{
		Code:        params.Code,
		Module:      params.Module,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.TransferInfo(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.TransferInfoResponse{
		ProductId:              resp.ProductId,
		CertificateNum:         resp.CertificateNum,
		Status:                 resp.Status,
		ErrorMessage:           resp.ErrorMessage,
		TransferCertificateNum: resp.TransferCertificateNum,
		TransferCertificateUrl: resp.TransferCertificateUrl,
	}, nil
}

func (r Rights) Revoke(ctx context.Context, params *dto.RevokeRequest) (*dto.RevokeResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "Revoke")

	req := rights.RevokeRequest{
		Code:           params.Code,
		Module:         params.Module,
		ProjectId:      params.ProjectID,
		OperationId:    params.OperationID,
		ProductId:      params.ProductID,
		CertificateNum: params.CertificateNum,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.Revoke(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.RevokeResponse{OperationID: resp.OperationId}, nil
}
func (r Rights) EditRevoke(ctx context.Context, params *dto.EditRevokeRequest) (*dto.EditRevokeResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "EditRevoke")

	req := rights.EditRevokeRequest{
		Code:        params.Code,
		Module:      params.Module,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.EditRevoke(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}
	return &dto.EditRevokeResponse{}, nil
}
func (r Rights) RevokeInfo(ctx context.Context, params *dto.RevokeInfoRequest) (*dto.RevokeInfoResponse, error) {
	logger := r.logger.WithField("params", params).WithField("func", "RevokeInfo")

	req := rights.RevokeInfoRequest{
		Code:        params.Code,
		Module:      params.Module,
		ProjectId:   params.ProjectID,
		OperationId: params.OperationID,
	}
	grpcClient, ok := initialize.RightsClientMap[constant.RightsMap[params.RegisterType]]
	if !ok {
		logger.Error(errors2.ErrService)
		return nil, errors2.New(errors2.ClientParams, "invalid register_type")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(constant.GrpcTimeout))
	defer cancel()
	resp, err := grpcClient.RevokeInfo(ctx, &req)
	if err != nil {
		logger.Error("grpc request failed")
		return nil, err
	}
	if resp == nil {
		return nil, errors2.New(errors2.InternalError, errors2.ErrGrpc)
	}

	return &dto.RevokeInfoResponse{
		ProductId:            resp.ProductId,
		CertificateNum:       resp.CertificateNum,
		Status:               resp.Status,
		ErrMessage:           resp.ErrorMessage,
		RevokeCertificateNum: resp.RevokeCertificateNum,
	}, nil
}
