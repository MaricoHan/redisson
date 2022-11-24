package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

func AuthData(ctx context.Context) (vo.AuthData, error) {
	authDataString := ctx.Value("X-Auth-Data")
	authDataSlice, ok := authDataString.([]string)
	if !ok {
		return vo.AuthData{}, fmt.Errorf("missing project parameters")
	}
	var authData vo.AuthData
	err := json.Unmarshal([]byte(authDataSlice[0]), &authData)
	if err != nil {
		return vo.AuthData{}, err
	}
	return authData, nil
}

func HeaderAuthData(header *http.Header) (vo.AuthData, error) {
	authDataString := header.Get("X-Auth-Data")
	var authData vo.AuthData
	err := json.Unmarshal([]byte(authDataString), &authData)
	if err != nil {
		return vo.AuthData{}, err
	}
	return authData, nil
}

func StrNameCheck(str string) bool {
	ok, err := regexp.MatchString("^[\u4E00-\u9FA5A-Za-z0-9]{1,20}$", str)
	if !ok || err != nil {
		return false
	}
	return true
}

func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.Trim(str, " \\t\\n\\r\\v\\f")
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}
