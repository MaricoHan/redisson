package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func StrNameCheck(str string) bool {
	ok, err := regexp.MatchString("^[\u4E00-\u9FA5A-Za-z0-9]{1,20}$", str)
	if !ok || err != nil {
		return false
	}
	return true
}

func Post(ctx context.Context, url, apiKey, apiSignature, customerID, timestamp string, body map[string]interface{}) (*http.Response, error) {
	bodys, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(bodys)))
	if err != nil {
		return nil, err
	}
	request.Header.Add("content-type", "application/json")
	request.Header.Add("X-Api-Key", apiKey)
	request.Header.Add("X-Signature", apiSignature)
	request.Header.Add("X-Timestamp", timestamp)
	request.Header.Add("X-Customer-ID", customerID)
	return http.DefaultClient.Do(request)
}

func Get(ctx context.Context, url, apiKey, apiSignature, customerID, timestamp string, body map[string]interface{}) (*http.Response, error) {
	var bodys []byte
	var err error
	var request *http.Request
	if body != nil {
		bodys, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequestWithContext(ctx, http.MethodGet, url, strings.NewReader(string(bodys)))
	} else {
		request, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	}
	if err != nil {
		return nil, err
	}
	request.Header.Add("content-type", "application/json")
	request.Header.Add("X-Api-Key", apiKey)
	request.Header.Add("X-Signature", apiSignature)
	request.Header.Add("X-Timestamp", timestamp)
	request.Header.Add("X-Customer-ID", customerID)
	return http.DefaultClient.Do(request)
}

func TimeToUnix(e time.Time) string {
	timeUnix, _ := time.Parse(constant.TimeLayout, e.Format(constant.TimeLayout))
	timeUnixString := strconv.FormatInt(timeUnix.UnixNano()/1e6, 10)
	return timeUnixString
}
