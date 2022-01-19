package handlers

import "context"

const timeLayout = "2006-01-02 15:04:05"

type base struct {
}

func (h base) AppID(ctx context.Context) uint64 {
	appID := ctx.Value("X-App-ID")
	return appID.(uint64)
}

type pageBasic struct {
}

func (h pageBasic) Offset(ctx context.Context) int64 {
	offset := ctx.Value("offset")
	if offset == nil {
		return 1
	}
	return offset.(int64)
}

func (h pageBasic) Limit(ctx context.Context) int64 {
	limit := ctx.Value("limit")
	if limit == nil {
		return 10
	}
	return limit.(int64)
}

func (h pageBasic) StartDate(ctx context.Context) string {
	endDate := ctx.Value("start_date")
	if endDate == nil {
		return ""
	}
	return endDate.(string)
}

func (h pageBasic) EndDate(ctx context.Context) string {
	endDate := ctx.Value("end_date")
	if endDate == nil {
		return ""
	}
	return endDate.(string)
}

func (h pageBasic) SortBy(ctx context.Context) string {
	sortBy := ctx.Value("sort_by")
	if sortBy == nil {
		return "DATE_DESC"
	}
	return sortBy.(string)
}
