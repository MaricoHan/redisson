package handlers

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"

	pb "gitlab.bianjie.ai/avata/chains/api/pb/account"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

type base struct {
}

func (h base) AuthData(ctx context.Context) vo.AuthData {
	authDataString := ctx.Value("X-Auth-Data")
	authDataSlice, ok := authDataString.([]string)
	if !ok {
		return vo.AuthData{}
	}
	var authData vo.AuthData
	err := json.Unmarshal([]byte(authDataSlice[0]), &authData)
	if err != nil {
		log.Error("auth data Error: ", err)
		return vo.AuthData{}
	}
	return authData
}

func (h base) UriCheck(uri string) error {
	if len([]rune(uri)) == 0 {
		return nil
	}
	if len([]rune(uri)) > 256 {
		return errors2.New(errors2.ClientParams, errors2.ErrURILen)
	}

	isUri := govalidator.IsRequestURI(uri)
	if !isUri {
		return errors2.New(errors2.ClientParams, errors2.ErrURIFormat)
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
		return 0, errors2.New(errors2.ClientParams, errors2.ErrOffset)
	}
	if offsetInt < 0 {
		return 0, errors2.New(errors2.ClientParams, errors2.ErrOffsetInt)
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
		return 10, errors2.New(errors2.ClientParams, errors2.ErrLimitParam)
	}
	if limitInt < 1 || limitInt > 50 {
		return 10, errors2.New(errors2.ClientParams, errors2.ErrLimitParamInt)
	}
	return limitInt, nil
}

func (h pageBasic) StartDate(ctx context.Context) string {
	startDate := ctx.Value("start_date")
	if startDate == "" || startDate == nil {
		return ""
	}
	return startDate.(string)
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
		return pb.SORT_name[0]
	}
	return sortBy.(string)
}
