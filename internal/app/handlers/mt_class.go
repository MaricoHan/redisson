package handlers

import (
	"context"
	mt2 "gitlab.bianjie.ai/avata/open-api/internal/app/models/dto/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/mt"
)

type IMtClass interface {
	List(ctx context.Context, _ interface{}) (interface{}, error)
	Show(ctx context.Context, _ interface{}) (interface{}, error)
}

type MtClass struct {
	base
	pageBasic
	svc mt.IMTClass
}

func NewMTClass(svc mt.IMTClass) *MtClass {
	return &MtClass{svc: svc}
}

// Classes return class list
func (h MtClass) List(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := mt2.MTClassListRequest{
		MtClassId:   h.MtClassId(ctx),
		MtClassName: h.MtClassName(ctx),
		Owner:       h.Owner(ctx),
		TxHash:      h.TxHash(ctx),
		ProjectID:   authData.ProjectId,
		ChainID:     authData.ChainId,
		PlatFormID:  authData.PlatformId,
		Module:      authData.Module,
		Code:        authData.Code,
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

	startDateR := h.StartDate(ctx)
	if startDateR != "" {
		params.StartDate = startDateR
	}

	endDateR := h.EndDate(ctx)
	if endDateR != "" {
		params.EndDate = endDateR
	}

	params.SortBy = h.SortBy(ctx)
	//switch h.SortBy(ctx) {
	//case "DATE_ASC":
	//	params.SortBy = "DATE_ASC"
	//case "DATE_DESC":
	//	params.SortBy = "DATE_DESC"
	//default:
	//	return nil, constant.NewAppError(constant.RootCodeSpace, constant.ClientParamsError, constant.ErrSortBy)
	//}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.List(&params)
}

// ClassByID return class
func (h MtClass) Show(ctx context.Context, _ interface{}) (interface{}, error) {
	// 校验参数 start
	authData := h.AuthData(ctx)
	params := mt2.MTClassShowRequest{
		ProjectID: authData.ProjectId,
		ClassID:   h.MtClassId(ctx),
		Status:    h.Status(ctx),
		Module:    authData.Module,
		Code:      authData.Code,
	}

	// 校验参数 end
	// 业务数据入库的地方
	return h.svc.Show(&params)
}

func (h MtClass) MtClassId(ctx context.Context) string {
	val := ctx.Value("mt_class_id")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MtClass) MtClassName(ctx context.Context) string {
	val := ctx.Value("mt_class_name")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MtClass) Owner(ctx context.Context) string {
	val := ctx.Value("owner")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MtClass) TxHash(ctx context.Context) string {
	val := ctx.Value("tx_hash")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MtClass) Timestamp(ctx context.Context) string {
	val := ctx.Value("timestamp")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MtClass) Status(ctx context.Context) string {
	val := ctx.Value("status")
	if val == nil {
		return ""
	}
	return val.(string)
}
func (h MtClass) MtCount(ctx context.Context) uint64 {
	val := ctx.Value("mt_count")
	if val == nil {
		return 0
	}
	return val.(uint64)
}
