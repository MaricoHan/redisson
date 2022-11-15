package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

func AuthData(ctx context.Context) vo.AuthData {
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

func StrNameCheck(str string) bool {
	ok, err := regexp.MatchString("^[\u4E00-\u9FA5A-Za-z0-9]{1,20}$", str)
	if !ok || err != nil {
		return false
	}
	return true
}

func Post(ctx context.Context, url, apiKey, apiSignature, customerID string, body map[string]interface{}) (*http.Response, error) {
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
	request.Header.Add("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	request.Header.Add("X-Customer-ID", customerID)
	return http.DefaultClient.Do(request)
}

func Get(ctx context.Context, url, apiKey, apiSignature, customerID string, body map[string]interface{}) (*http.Response, error) {
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
	request.Header.Add("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	request.Header.Add("X-Customer-ID", customerID)
	return http.DefaultClient.Do(request)
}
