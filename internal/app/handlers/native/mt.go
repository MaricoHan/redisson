package native

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/base"

	dto "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/native/mt"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/native"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
)

type IMT interface {
	Issue(ctx context.Context, request interface{}) (response interface{}, err error)
	Mint(ctx context.Context, request interface{}) (response interface{}, err error)
	BatchMint(ctx context.Context, request interface{}) (response interface{}, err error)

	Edit(ctx context.Context, request interface{}) (response interface{}, err error)
	Burn(ctx context.Context, request interface{}) (response interface{}, err error)
	Transfer(ctx context.Context, request interface{}) (response interface{}, err error)

	Show(ctx context.Context, request interface{}) (response interface{}, err error)
	List(ctx context.Context, request interface{}) (response interface{}, err error)
	Balances(ctx context.Context, request interface{}) (response interface{}, err error)
}
type MT struct {
	base.Base
	base.PageBasic
	svc native.IMT
}

func NewMT(svc native.IMT) *MT {
	return &MT{svc: svc}
}

func (h MT) Issue(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.IssueRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
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
		OperationID: req.OperationID,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.Issue(ctx, &param)
}

func (h MT) Mint(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.MintRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
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
		Amount:      req.Amount,
		Recipient:   req.Recipient,
		OperationID: req.OperationID,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.Mint(ctx, &param)
}

func (h MT) BatchMint(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.BatchMintRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.BatchMintRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		ClassID:     h.ClassID(ctx),
		MTID:        h.MTID(ctx),
		Recipients:  req.Recipients,
		OperationID: req.OperationID,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.BatchMint(ctx, &param)
}

func (h MT) Edit(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.EditRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.EditRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		Owner:       h.Owner(ctx),
		ClassId:     h.ClassID(ctx),
		MTID:        h.MTID(ctx),
		Data:        req.Data,
		OperationID: req.OperationID,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.Edit(ctx, &param)
}

func (h MT) Burn(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.BurnRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.BurnRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		Owner:       h.Owner(ctx),
		ClassID:     h.ClassID(ctx),
		MtID:        h.MTID(ctx),
		Amount:      req.Amount,
		OperationID: req.OperationID,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.Burn(ctx, &param)
}

func (h MT) BatchBurn(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.BatchBurnRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.BatchBurnRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		Owner:       h.Owner(ctx),
		Mts:         req.Mts,
		OperationID: req.OperationID,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.BatchBurn(ctx, &param)
}

func (h MT) Transfer(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.TransferRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.MTTransferRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		Owner:       h.Owner(ctx),
		ClassId:     h.ClassID(ctx),
		MtId:        h.MTID(ctx),
		Amount:      req.Amount,
		Recipient:   req.Recipient,
		OperationID: req.OperationID,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.Transfer(ctx, &param)
}

func (h MT) BatchTransfer(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 接收请求
	req, ok := request.(*vo.BatchTransferRequest)
	if !ok {
		log.Debugf("failed to assert : %v", request)
		return nil, errors2.New(errors2.ClientParams, errors2.ErrClientParams)
	}
	req.OperationID = strings.TrimSpace(req.OperationID)

	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.MTBatchTransferRequest{
		Code:        authData.Code,
		Module:      authData.Module,
		ProjectID:   authData.ProjectId,
		Owner:       h.Owner(ctx),
		Mts:         req.Mts,
		OperationID: req.OperationID,
		AccessMode:  authData.AccessMode,
	}

	return h.svc.BatchTransfer(ctx, &param)
}

func (h MT) Show(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := h.AuthData(ctx)
	param := dto.MTShowRequest{
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		ClassID:    h.ClassID(ctx),
		MTID:       h.MTID(ctx),
		AccessMode: authData.AccessMode,
	}

	return h.svc.Show(ctx, &param)
}

func (h MT) List(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := h.AuthData(ctx)
	params := dto.MTListRequest{
		ProjectID:  authData.ProjectId,
		MtId:       h.ID(ctx),
		ClassId:    h.ClassID(ctx),
		Issuer:     h.Issuer(ctx),
		TxHash:     h.TxHash(ctx),
		Module:     authData.Module,
		Code:       authData.Code,
		AccessMode: authData.AccessMode,
	}
	params.PageKey = h.PageKey(ctx)
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
	return h.svc.List(ctx, &params)
}

func (h MT) Balances(ctx context.Context, request interface{}) (response interface{}, err error) {
	// 获取账户基本信息
	authData := h.AuthData(ctx)
	params := dto.MTBalancesRequest{
		ProjectID:  authData.ProjectId,
		Module:     authData.Module,
		Code:       authData.Code,
		MtId:       h.ID(ctx),
		ClassId:    h.ClassID(ctx),
		Account:    h.Account(ctx),
		AccessMode: authData.AccessMode,
	}
	params.PageKey = h.PageKey(ctx)

	limit, err := h.Limit(ctx)
	if err != nil {
		return nil, err
	}
	params.Limit = limit

	if params.Limit == 0 {
		params.Limit = 10
	}

	return h.svc.Balances(ctx, &params)
}

func (MT) ClassID(ctx context.Context) string {
	classId := ctx.Value("class_id")

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
func (MT) ID(ctx context.Context) string {
	mtID := ctx.Value("id")

	if mtID == nil {
		return ""
	}
	return mtID.(string)
}
func (MT) Account(ctx context.Context) string {
	account := ctx.Value("account")

	if account == nil {
		return ""
	}
	return account.(string)
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
