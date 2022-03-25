package mw

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"gitlab.bianjie.ai/irita-paas/open-api/config"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/irita-paas/orms/orm-nft/models"
)

// 误差时间
const timeInterval = 30

// AuthMiddleware recover the panic error
func AuthMiddleware(h http.Handler) http.Handler {
	return authHandler{h}
}

type authHandler struct {
	next http.Handler
}

func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("user request", "method:", r.Method, "url:", r.URL.Path)

	createTime := time.Now()
	defer func(createTime time.Time) {
		//监控响应时间
		interval := time.Now().Sub(createTime)
		metric.NewPrometheus().ApiHttpRequestRtSeconds.With([]string{
			"method",
			r.Method,
			"uri",
			r.RequestURI,
		}...).Observe(float64(interval))
	}(createTime)

	appKey := r.Header.Get("X-Api-Key")
	projectKeyResult, err := models.TProjectKeys(
		qm.Select(models.TProjectKeyColumns.APISecret),
		qm.Select(models.TProjectKeyColumns.ProjectID),
		models.TProjectKeyWhere.APIKey.EQ(appKey),
	).OneG(context.Background())
	if err != nil {
		log.Error("server http", "project keys error:", err)
		writeForbiddenResp(w)
		return
	}

	// 获取project
	project, err := models.TProjects(
		models.TProjectWhere.ID.EQ(projectKeyResult.ProjectID)).OneG(context.Background())
	if err != nil {
		log.Error("server http", "project error:", err)
		writeForbiddenResp(w)
		return
	}

	// 获取链信息
	chain, err := models.TChains(models.TChainWhere.ID.EQ(project.ChainID)).OneG(context.Background())
	if err != nil {
		log.Error("server http", "chain error:", err)
		writeForbiddenResp(w)
		return
	}

	authData := AuthData{
		ProjectId:  project.ID,
		ChainId:    chain.ID,
		PlatformId: uint64(project.PlatformID.Int64),
		Module:     chain.Module,
	}
	authDataBytes, _ := json.Marshal(authData)
	// 1. 获取 header 中的时间戳
	reqTimestampStr := r.Header.Get("X-Timestamp")

	//// 1.1 判断时间误差
	// todo
	//生产的时候打开
	//reqTimestampInt, err := strconv.ParseInt(reqTimestampStr, 10, 64)
	//if err != nil {
	//	writeBadRequestResp(w, types.ErrParams)
	//	return
	//}
	//
	//curTimestamp := time.Now().Unix()
	//if curTimestamp-reqTimestampInt > timeInterval || curTimestamp < reqTimestampInt {
	//	writeBadRequestResp(w, types.ErrParams)
	//	return
	//}

	reqSignature := r.Header.Get("X-Signature")
	// 2. 验证签名
	// todo
	// 生产的时候打开
	if config.Get().Server.SignatureAuth && !h.Signature(r, projectKeyResult.APISecret, reqTimestampStr, reqSignature) {
		writeForbiddenResp(w)
		return
	}
	log.Info("signature: ", h.Signature(r, projectKeyResult.APISecret, reqTimestampStr, reqSignature))
	r.Header.Set("X-App-Id", fmt.Sprintf("%d", authData.ProjectId))
	r.Header.Set("X-Auth-Data", fmt.Sprintf("%s", string(authDataBytes)))
	h.next.ServeHTTP(w, r)
}

func (h authHandler) Signature(r *http.Request, apiSecret string, timestamp string, signature string) bool {

	// 获取 path params
	params := map[string]interface{}{}
	//for k, v := range mux.Vars(r) {
	//	params[k] = v
	//}
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
