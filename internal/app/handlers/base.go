package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"

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

func (h base) ValidateTag(tags map[string]interface{}) ([]byte, error) {
	var tagBytes []byte
	if len(tags) > 0 {
		tagBytes, _ = json.Marshal(tags)
		tag := string(tagBytes)
		if _, err := h.IsValTag(tag); err != nil {
			return tagBytes, errors2.New(errors2.ClientParams, fmt.Sprintf("invalid tag :%s", err.Error()))
			//return tagBytes, constant.NewAppError(constant.RootCodeSpace, constant.ClientParamsError, err.Error())
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
		if len([]rune(sValue)) == 0 || len([]rune(sValue)) > 64 {
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
		return "DATE_DESC"
	}
	return sortBy.(string)
}
