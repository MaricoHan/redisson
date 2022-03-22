package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"unicode"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/mw"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"

	"github.com/asaskevich/govalidator"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

const (
	timeLayout           = "2006-01-02 15:04:05"
	timeLayoutWithoutHMS = "2006-01-02"
	SqlNoFound           = "records not exist"
)

type base struct {
}

func (h base) AuthData(ctx context.Context) mw.AuthData {
	authDataString := ctx.Value("X-Auth-Data")
	authDataSlice, ok := authDataString.([]string)
	if !ok {
		return mw.AuthData{}
	}
	var authData mw.AuthData
	err := json.Unmarshal([]byte(authDataSlice[0]), &authData)
	if err != nil {
		log.Error("auth data Error: ", err)
		return mw.AuthData{}
	}
	return authData
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

func (h base) ValidateTag(tags map[string]interface{}) ([]byte, error) {
	var tagBytes []byte
	if len(tags) > 0 {
		tagBytes, _ = json.Marshal(tags)
		tag := string(tagBytes)
		if _, err := h.IsValTag(tag); err != nil {
			return tagBytes, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, err.Error())
		}
	}
	return tagBytes, nil
}

func (h base) IsValTag(tag string) (bool, error) {
	//校验tag是否是json格式
	if tag[0] != '{' || !json.Valid([]byte(tag)) {
		return false, errors.New("invalid json format")
	}
	//解析tag为map
	var f interface{}
	json.Unmarshal([]byte(tag), &f)
	tagMap := f.(map[string]interface{})
	if len(tagMap) > 3 {
		return false, errors.New("at most 3 key-value in a tag at a time")
	}
	//校验tagMap的各个key-value
	for key, value := range tagMap {
		sValue, ok := value.(string)
		if !ok {
			return false, errors.New("value must be string")
		}
		sValue = strings.TrimSpace(sValue)
		key = strings.TrimSpace(key)
		if len([]rune(key)) < 6 || len([]rune(key)) > 12 {
			return false, errors.New("key’s length must between 6 and 12")
		}
		for _, s := range key {
			if (s < '0' || s > '9') && (s < 'A' || s > 'Z') && (s < 'a' || s > 'z') && !unicode.Is(unicode.Han, s) {
				return false, errors.New("key must contain only letters , numbers and Chinese characters")
			}
		}
		if len(sValue) == 0 || len(sValue) > 64 {
			return false, errors.New("value’s length must between 1 and 64")
		}
		for _, s := range sValue {
			if (s < '0' || s > '9') && (s < 'A' || s > 'Z') && (s < 'a' || s > 'z') {
				return false, errors.New("value must contain only letters and numbers")
			}
		}
	}
	return true, nil
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
