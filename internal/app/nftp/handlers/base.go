package handlers

import (
	"context"
	"encoding/json"
	"github.com/friendsofgo/errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/asaskevich/govalidator"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

const timeLayout = "2006-01-02 15:04:05"
const timeLayoutWithoutHMS = "2006-01-02"

type base struct {
}

func (h base) ChainID(ctx context.Context) uint64 {
	keysList := ctx.Value("X-App-Id")
	keysListString, ok := keysList.([]string)
	if !ok {
		return 0
	}
	ChainID, _ := strconv.ParseInt(keysListString[0], 10, 64)
	return uint64(ChainID)
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
func (h base)IsValTag(tag string) (bool, error) {
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
		sValue = strings.TrimSpace(sValue)
		key = strings.TrimSpace(key)
		if !ok {
			return false, errors.New("value must be string")
		}
		if len(key) > 12 || len(key) < 6 {
			return false, errors.New("key’s length must between 6 and 12")
		}
		for _, s := range key {
			if (s < '0' || s > '9') && (s < 'A' || s > 'Z') && (s < 'a' || s > 'z') && !unicode.Is(unicode.Han, s) {
				return false, errors.New("key must contain only letters , numbers and Chinese characters")
			}
		}
		if len(sValue) > 64 {
			return false, errors.New("The maximum length of value is 64")
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
