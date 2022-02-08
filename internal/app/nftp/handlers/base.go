package handlers

import (
	"context"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"strconv"
)

const timeLayout = "2006-01-02 15:04:05"
const timeLayoutWithoutHMS = "2006-01-02"

type base struct {
}

func (h base) AppID(ctx context.Context) uint64 {
	keysList := ctx.Value("X-App-Id")
	keysListString, ok := keysList.([]string)
	if !ok {
		return 0
	}
	appID, _ := strconv.ParseInt(keysListString[0], 10, 64)
	return uint64(appID)
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
		return 0, err
	}
	if offsetInt < 0 {
		return 0, types.ErrParams
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
		return 10, err
	}
	if limitInt < 1 || limitInt > 50 {
		return 10, types.ErrParams
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
