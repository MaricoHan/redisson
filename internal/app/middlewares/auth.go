package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/chain"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/project"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// 误差时间
const timeInterval = 3000

// AuthMiddleware recover the panic error
func AuthMiddleware(h http.Handler) http.Handler {
	return authHandler{h}
}

type authHandler struct {
	next http.Handler
}

func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("user request", " method:", r.Method, " url:", r.URL.Path)

	//createTime := time.Now()
	//defer func(createTime time.Time) {
	//	//监控响应时间
	//	interval := time.Now().Sub(createTime)
	//	metric.NewPrometheus().ApiHttpRequestRtSeconds.With([]string{
	//		"method",
	//		r.Method,
	//		"uri",
	//		r.RequestURI,
	//	}...).Observe(float64(interval))
	//}(createTime)

	appKey := r.Header.Get("X-Api-Key")

	//查询缓存
	var projectInfo entity.Project
	err := initialize.RedisClient.GetObject(fmt.Sprintf("%s%s", constant.KeyProjectApikey, appKey), &projectInfo)
	if err != nil {
		log.Error("server http", "get project cache error:", err)
		writeInternalResp(w)
		return
	}
	if projectInfo.ID < 1 {
		//查询project信息
		projectRepo := project.NewProjectRepo(initialize.MysqlDB)
		projectInfo, err = projectRepo.GetProjectByApiKey(appKey)
		if err != nil {
			log.Error("server http", "project error:", err)
			writeInternalResp(w)
			return
		}

		if projectInfo.ID > 0 {
			// save cache
			if err := initialize.RedisClient.SetObject(fmt.Sprintf("%s%s", constant.KeyProjectApikey, appKey), projectInfo, time.Hour*24); err != nil {
				log.Error("server http", "save project cache error:", err)
				writeInternalResp(w)
				return
			}
		}

	}

	if projectInfo.ID == 0 {
		log.Error("server http:", constant.ErrApikey)
		writeForbiddenResp(w, constant.ErrApikey)
		return
	}

	//查询缓存
	var chainInfo entity.Chain
	err = initialize.RedisClient.GetObject(fmt.Sprintf("%s%d", constant.KeyChain, projectInfo.ChainID), &chainInfo)
	if err != nil {
		log.Error("server http", "get chain cache error:", err)
		writeInternalResp(w)
		return
	}
	if chainInfo.ID < 1 {
		// 获取链信息
		chainRepo := chain.NewChainRepo(initialize.MysqlDB)
		chainInfo, err = chainRepo.QueryChainById(projectInfo.ChainID)
		if err != nil {
			log.Error("server http", "chain error:", err)
			writeInternalResp(w)
			return
		}
		if chainInfo.ID > 0 {
			// save cache
			if err := initialize.RedisClient.SetObject(fmt.Sprintf("%s%d", constant.KeyChain, projectInfo.ChainID), chainInfo, time.Hour*24); err != nil {
				log.Error("server http", "save project cache error:", err)
				writeInternalResp(w)
				return
			}
		}
	}

	if chainInfo.ID == 0 {
		log.Error("server http", "project not exist:")
		writeInternalResp(w)
		return
	}

	authData := vo.AuthData{
		ProjectId:  uint64(projectInfo.ID),
		ChainId:    uint64(chainInfo.ID),
		PlatformId: projectInfo.UserID,
		Module:     chainInfo.Module,
		Code:       chainInfo.Code,
	}

	authDataBytes, _ := json.Marshal(authData)
	// 1. 判断时间误差
	if configs.Cfg.App.TimestampAuth {
		reqTimestampStr := r.Header.Get("X-Timestamp")
		reqTimestampInt, err := strconv.ParseInt(reqTimestampStr, 10, 64)
		if err != nil {
			writeBadRequestResp(w, constant.ErrParams)
			return
		}

		curTimestamp := time.Now().Unix()
		if curTimestamp-reqTimestampInt > timeInterval || curTimestamp < reqTimestampInt {
			writeBadRequestResp(w, constant.ErrTimestamp)
			return
		}
	}

	reqSignature := r.Header.Get("X-Signature")
	// 2. 验证签名
	if configs.Cfg.App.SignatureAuth {
		reqTimestampStr := r.Header.Get("X-Timestamp")
		if !h.Signature(r, projectInfo.ApiSecret, reqTimestampStr, reqSignature) {

			writeForbiddenResp(w, "")
			return
		}
		log.Debugf("signature: %v", h.Signature(r, projectInfo.ApiSecret, reqTimestampStr, reqSignature))
	}

	r.Header.Set("X-App-Id", fmt.Sprintf("%d", authData.ProjectId))
	r.Header.Set("X-Auth-Data", fmt.Sprintf("%s", string(authDataBytes)))
	h.next.ServeHTTP(w, r)
}

func (h authHandler) Signature(r *http.Request, apiSecret string, timestamp string, signature string) bool {

	// 获取 path params
	params := map[string]interface{}{}
	params["path_url"] = r.URL.Path

	// 获取 query params
	for k, v := range r.URL.Query() {
		k = "query_" + k
		params[k] = v[0]
	}

	// 获取 body params
	// 把request的内容读取出来
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	// 把刚刚读出来的再写进去
	if bodyBytes != nil {
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	paramsBody := map[string]interface{}{}
	_ = json.Unmarshal(bodyBytes, &paramsBody)
	hexHash := hash(timestamp + apiSecret)

	for k, v := range paramsBody {
		k = "body_" + k
		params[k] = v
	}
	// sort params
	//sortParams := sortMapParams(params)
	sortParams := params
	if sortParams != nil {
		sortParamsBytes, _ := json.Marshal(sortParams)
		hexHash = hash(string(sortParamsBytes) + timestamp + apiSecret)
	}
	if hexHash != signature {
		return false
	}
	return true

}

func hash(oriText string) string {
	oriTextHashBytes := sha256.Sum256([]byte(oriText))
	return hex.EncodeToString(oriTextHashBytes[:])
}

func sortMapParams(params map[string]interface{}) map[string]interface{} {
	keys := make([]string, len(params))
	i := 0
	for k, _ := range params {
		keys[i] = k
		i++
	}
	if i == 0 {
		return nil
	}
	sort.Strings(keys)
	sortMap := map[string]interface{}{}
	for _, k := range keys {
		sortMap[k] = params[k]
	}
	return sortMap
}
