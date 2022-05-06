package utils

import (
	"context"
	"encoding/json"
	"regexp"

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
