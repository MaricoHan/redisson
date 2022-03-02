package handlers

import (
	"context"
	"strconv"

	"github.com/asaskevich/govalidator"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

const timeLayout = "2006-01-02 15:04:05"
const timeLayoutWithoutHMS = "2006-01-02"

type base struct {
}

func (h base) ChainID(ctx context.Context) uint64 {
	rec := ctx.Value("chain_id")
	chainId, ok := rec.(string)
	if !ok {
		return 0
	}
	res, _ := strconv.ParseInt(chainId, 10, 64)
	return uint64(res)
}

func (h base) UriCheck(uri string) error {
	if len([]rune(uri)) == 0 {
		return nil
	}
	if len([]rune(uri)) > 256 {
		return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrUriLen)
	}

	isUri := govalidator.IsRequestURI(uri)
	if !isUri {
		return types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrUri)
	}

	return nil
}

type pageBasic struct {
}

func (h pageBasic) Offset(ctx context.Context) (int64, error) {
	offset := ctx.Value("offset")
	if offset == "" || offset == nil {
		return 0, nil
	}
	offsetInt, err := strconv.ParseInt(offset.(string), 10, 64)
	if err != nil {
		return 0, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOffset)
	}
	if offsetInt < 0 {
		return 0, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrOffsetInt)
	}
	return offsetInt, nil
}

func (h pageBasic) Limit(ctx context.Context) (int64, error) {
	limit := ctx.Value("limit")
	if limit == "" || limit == nil {
		return 10, nil
	}
	limitInt, err := strconv.ParseInt(limit.(string), 10, 64)
	if err != nil {
		return 10, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrLimitParam)
	}
	if limitInt < 1 || limitInt > 50 {
		return 10, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrLimitParamInt)
	}
	return limitInt, nil
}

func (h pageBasic) StartDate(ctx context.Context) string {
	endDate := ctx.Value("start_date")
	if endDate == "" || endDate == nil {
		return ""
	}
	return endDate.(string)
}

func (h pageBasic) EndDate(ctx context.Context) string {
	endDate := ctx.Value("end_date")
	if endDate == "" || endDate == nil {
		return ""
	}
	return endDate.(string)
}

func (h pageBasic) SortBy(ctx context.Context) string {
	sortBy := ctx.Value("sort_by")
	if sortBy == "" || sortBy == nil {
		return "DATE_DESC"
	}
	return sortBy.(string)
}
