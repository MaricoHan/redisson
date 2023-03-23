package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/project"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/cache"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/configs"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/avata/utils/commons/aes"
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
	appKey := r.Header.Get("X-Api-Key")

	log := initialize.Log.WithFields(map[string]interface{}{
		"function": "ServeHTTP",
		"method":   r.Method,
		"url":      r.URL.Path,
		"api-key":  appKey,
	})

	// 查询缓存
	projectInfo, existWalletService, err := cache.NewCache().Project(appKey)
	if err != nil {
		log.WithError(err).Error("project from cache")
		writeInternalResp(w)
		return
	}

	if projectInfo.Id == 0 {
		log.Error(constant.ErrApikey)
		writeForbiddenResp(w, constant.ErrApikey)
		return
	}

	key := fmt.Sprintf("balance:%d", projectInfo.Id)
	exists, err := initialize.RedisClient.Exists(key)
	if err != nil {
		log.WithError(err).Error("redis exists")
		writeInternalResp(w)
		return
	}
	if exists {
		writeNotEnoughAmount(w)
		return
	}

	// 查询缓存
	chainInfo, err := cache.NewCache().Chain(projectInfo.ChainId)
	if err != nil {
		log.WithError(err).Error("chain from cache")
		writeInternalResp(w)
		return
	}

	if chainInfo.Id == 0 {
		log.Error("project not exist")
		writeInternalResp(w)
		return
	}

	authData := vo.AuthData{
		ProjectId:  uint64(projectInfo.Id),
		ChainId:    uint64(chainInfo.Id),
		PlatformId: uint64(projectInfo.UserId),
		Module:     chainInfo.Module,
		Code:       chainInfo.Code,
		AccessMode: projectInfo.AccessMode,
		UserId:     uint64(projectInfo.UserId),
	}

	// 判断项目参数版本号
	if projectInfo.Version == entity.Version1 {
		log.Error("project version not implemented")
		writeNotFoundRequestResp(w, constant.ErrUnSupported)
		return
	} else if projectInfo.Version == entity.VersionStage {
		authData.Code = constant.IritaOPB
		authData.Module = constant.Native
	}

	// 如果项目关联钱包服务,转发到钱包微服务
	if existWalletService {
		authData.Code = constant.Wallet
		authData.Module = constant.Server
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
		secret, err := aes.Decode(projectInfo.ApiSecret, configs.Cfg.Project.SecretPwd)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"secret":   projectInfo.ApiSecret,
				"password": configs.Cfg.Project.SecretPwd,
			}).WithError(err).Error("decrypt api-secret failed")
			writeInternalResp(w)
			return
		}
		if !h.Signature(r, secret, reqTimestampStr, reqSignature) {
			metric.NewPrometheus().ApiServiceRequests.With([]string{
				"name", "open-api",
				"method", "/api.project.v1beta1.Project/Auth",
				"status", "401",
			}...).Add(1)
			writeForbiddenResp(w, "")
			return
		}
	}

	if strings.ContainsAny("/ns/", r.URL.Path) {
		// 域名请求，验证权限
		projectRepo := project.NewProjectRepo(initialize.MysqlDB)
		auth, err := projectRepo.ExistServices(projectInfo.Id, entity.ServiceTypeNS)
		if err != nil {
			log.WithError(err).Error("query service")
			writeInternalResp(w)
			return
		}
		if !auth {
			writeForbiddenResp(w, constant.AuthenticationFailed)
			return
		}
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
		bodyBytes, _ = io.ReadAll(r.Body)
	}
	// 把刚刚读出来的再写进去
	if bodyBytes != nil {
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	paramsBody := map[string]interface{}{}
	_ = json.Unmarshal(bodyBytes, &paramsBody)

	for k, v := range paramsBody {
		k = "body_" + k
		params[k] = v
	}

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(params)
	hexHash := hash(strings.TrimRight(bf.String(), "\n") + timestamp + apiSecret)
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
