package middlewares

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/chain"
	"gitlab.bianjie.ai/avata/open-api/internal/app/repository/db/project"
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

	// createTime := time.Now()
	// defer func(createTime time.Time) {
	//	//监控响应时间
	//	interval := time.Now().Sub(createTime)
	//	metric.NewPrometheus().ApiHttpRequestRtSeconds.With([]string{
	//		"method",
	//		r.Method,
	//		"uri",
	//		r.RequestURI,
	//	}...).Observe(float64(interval))
	// }(createTime)

	// 查询缓存
	var projectInfo entity.Project
	err := initialize.RedisClient.GetObject(fmt.Sprintf("%s%s", constant.KeyProjectApikey, appKey), &projectInfo)
	if err != nil {
		log.WithError(err).Error("get project from cache ")
		writeInternalResp(w)
		return
	}
	if projectInfo.Id < 1 {
		// 查询project信息
		projectRepo := project.NewProjectRepo(initialize.MysqlDB)
		projectInfo, err = projectRepo.GetProjectByApiKey(appKey)
		if err != nil {
			log.WithError(err).Error("get project from cache")
			writeInternalResp(w)
			return
		}

		if projectInfo.Id > 0 {
			// save cache
			if err := initialize.RedisClient.SetObject(fmt.Sprintf("%s%s", constant.KeyProjectApikey, appKey), projectInfo, time.Minute*5); err != nil {
				log.WithError(err).Error("save project cache")
				writeInternalResp(w)
				return
			}
		}

	}

	if projectInfo.Id == 0 {
		log.Error(constant.ErrApikey)
		writeForbiddenResp(w, constant.ErrApikey)
		return
	}

	// 查询缓存
	var chainInfo entity.Chain
	err = initialize.RedisClient.GetObject(fmt.Sprintf("%s%d", constant.KeyChain, projectInfo.ChainId), &chainInfo)
	if err != nil {
		log.WithError(err).Error("get chain from cache")
		writeInternalResp(w)
		return
	}
	if chainInfo.Id < 1 {
		// 获取链信息
		chainRepo := chain.NewChainRepo(initialize.MysqlDB)
		chainInfo, err = chainRepo.QueryChainById(uint64(projectInfo.ChainId))
		if err != nil {
			log.WithError(err).Error("query chain from db by id")
			writeInternalResp(w)
			return
		}
		if chainInfo.Id > 0 {
			// save cache
			if err := initialize.RedisClient.SetObject(fmt.Sprintf("%s%d", constant.KeyChain, projectInfo.ChainId), chainInfo, time.Minute*5); err != nil {
				log.WithError(err).Error("save project to cache")
				writeInternalResp(w)
				return
			}
		}
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
	}

	// DDC 不支持 NFT-批量、orders-批量、MT、版权服务
	if fmt.Sprintf("%s-%s", chainInfo.Code, chainInfo.Module) == constant.WenchangDDC {
		if strings.Contains(r.RequestURI, "/mt/") || strings.Contains(r.RequestURI, "/nft/batch/") || strings.Contains(r.RequestURI, "/orders/batch") || strings.Contains(r.RequestURI, "/rights/") {
			writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
			return
		}
	} else { // native
		// 非托管模式
		if projectInfo.AccessMode == entity.UNMANAGED {
			if fmt.Sprintf("%s-%s", chainInfo.Code, chainInfo.Module) == constant.IritaOPBNative {
				// 文昌链-天舟除 orders 都不支持
				if !strings.Contains(r.RequestURI, "/orders") {
					writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
					return
				}
			} else {
				// 文昌链-天和都不支持
				writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
				return
			}

		} else {
			// 托管不支持 orders
			if strings.Contains(r.RequestURI, "/orders") {
				writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
				return
			}
		}
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
	// sortParams := sortMapParams(params)
	sortParams := params
	if sortParams != nil {
		bf := bytes.NewBuffer([]byte{})
		jsonEncoder := json.NewEncoder(bf)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(sortParams)

		hexHash = hash(strings.TrimRight(bf.String(), "\n") + timestamp + apiSecret)
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
