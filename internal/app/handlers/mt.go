package handlers

import (
	"context"
	"encoding/json"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
	service "gitlab.bianjie.ai/avata/open-api/internal/app/services"
	"strings"

	log "github.com/sirupsen/logrus"

	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/mt"

	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IMT interface {
	Issue(ctx context.Context, request interface{}) (response interface{}, err error)
	Mint(ctx context.Context, request interface{}) (response interface{}, err error)
	Edit(ctx context.Context, request interface{}) (response interface{}, err error)
	Burn(ctx context.Context, request interface{}) (response interface{}, err error)
	Transfer(ctx context.Context, request interface{}) (response interface{}, err error)

	Show(ctx context.Context, request interface{}) (response interface{}, err error)
	List(ctx context.Context, request interface{}) (response interface{}, err error)
}
type MT struct {
	base
	pageBasic
	svc service.IMT
}

func NewMT(svc service.IMT) *MT {
	return &MT{svc: svc}
}

func (h MT) Issue(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.IssueRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 转换tag
	var tagBz []byte
	if len(req.Tag) > 0 {
		tagBz, _ = json.Marshal(req.Tag)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.IssueRequest{
		ProjectID:   authData.ProjectId,
		Module:      authData.Module,
		Code:        authData.Code,
		ClassID:     h.ClassID(ctx),
		Metadata:    req.Metadata,
		Amount:      req.Amount,
		Recipient:   req.Recipient,
		Tag:         string(tagBz),
		OperationID: req.OperationID,
	}

	return h.svc.Issue(&param)
}

func (h MT) Mint(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.MintRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 转换tag
	var tagBz []byte
	if len(req.Tag) > 0 {
		tagBz, _ = json.Marshal(req.Tag)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.MintRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		ClassID:     h.ClassID(ctx),
		MTID:        h.MTID(ctx),
		Recipients:  req.Recipients,
		Tag:         string(tagBz),
		OperationID: req.OperationID,
	}

	return h.svc.Mint(&param)
}

func (h MT) Edit(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.EditRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 转换tag
	var tagBz []byte
	if len(req.Tag) > 0 {
		tagBz, _ = json.Marshal(req.Tag)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.EditRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		Owner:       h.Owner(ctx),
		Mts:         req.Mts,
		Tag:         string(tagBz),
		OperationID: req.OperationID,
	}

	return h.svc.Edit(&param)
}

func (h MT) Burn(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.BurnRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 转换tag
	var tagBz []byte
	if len(req.Tag) > 0 {
		tagBz, _ = json.Marshal(req.Tag)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.BurnRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		Owner:       h.Owner(ctx),
		Mts:         req.Mts,
		Tag:         string(tagBz),
		OperationID: req.OperationID,
	}

	return h.svc.Burn(&param)
}

func (h MT) Transfer(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.TransferRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	// 转换tag
	var tagBz []byte
	if len(req.Tag) > 0 {
		tagBz, _ = json.Marshal(req.Tag)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.TransferRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		Owner:       h.Owner(ctx),
		Mts:         req.Mts,
		Tag:         string(tagBz),
		OperationID: req.OperationID,
	}

	return h.svc.Transfer(&param)
}

func (h MT) Show(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.MTShowRequest{
		ProjectID: authData.ProjectId,
		ClassID:   h.ClassID(ctx),
		MTID:      h.MTID(ctx),
		Module:    authData.Module,
		Code:      authData.Code,
	}

	return h.svc.Show(&param)
}

func (h MT) List(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := h.AuthData(ctx)
	params := dto.MTListRequest{
		ProjectID: authData.ProjectId,
		MtId:      h.MTID(ctx),
		MtClassId: h.ClassID(ctx),
		Issuer:    h.Issuer(ctx),
		TxHash:    h.TxHash(ctx),
		Module:    authData.Module,
		Code:      authData.Code,
	}
	offset, err := h.Offset(ctx)
	if err != nil {
		return nil, err
	}
	params.Offset = offset

	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, err
	}
	params.Limit = limit

	if params.Limit == 0 {
		params.Limit = 10
	}

	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		params.StartDate = startDateR
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {

		params.EndDate = endDateR
	}

	params.SortBy = h.SortBy(ctx)
	return h.svc.List(&params)
}

func (MT) ClassID(ctx context.Context) string {
	classId := ctx.Value("mt_class_id")

	if classId == nil {
		return ""
	}
	return classId.(string)
}
func (MT) MTID(ctx context.Context) string {
	mtID := ctx.Value("mt_id")

	if mtID == nil {
		return ""
	}
	return mtID.(string)
}

func (MT) Owner(ctx context.Context) string {
	val := ctx.Value("owner")

	if val == nil {
		return ""
	}
	return val.(string)
}

func (MT) Issuer(ctx context.Context) string {
	val := ctx.Value("issuer")

	if val == nil {
		return ""
	}
	return val.(string)
}
func (MT) TxHash(ctx context.Context) string {
	val := ctx.Value("tx_hash")

	if val == nil {
		return ""
	}
	return val.(string)
}
