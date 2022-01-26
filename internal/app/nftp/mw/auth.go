package mw

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gorilla/mux"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	appKey := r.Header.Get("X-Api-Key")
	appKeyResult, err := models.TAppKeys(
		qm.Select(models.TAppKeyColumns.APIKey),
		qm.Select(models.TAppKeyColumns.AppID),
		models.TAppKeyWhere.APIKey.EQ(appKey),
	).OneG(context.Background())
	if err != nil {
		writeForbiddenResp(w)
		return
	}
	// 1. 获取 header 中的时间戳
	reqTimestampStr := r.Header.Get("X-Timestamp")

	//// 1.1 判断时间误差
	//reqTimestampInt, err := strconv.ParseInt(reqTimestampStr, 10, 64)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//curTimestamp := time.Now().Unix()
	//if curTimestamp-reqTimestampInt > timeInterval || curTimestamp < reqTimestampInt {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	reqSignature := r.Header.Get("X-Signature")
	// 2. 验证签名
	// todo
	fmt.Println(h.Signature(r, appKeyResult.APIKey, reqTimestampStr, reqSignature))
	r.Header.Set("X-App-Id", fmt.Sprintf("%d", appKeyResult.AppID))
	h.next.ServeHTTP(w, r)
}

func (h authHandler) Signature(r *http.Request, appKey string, timestamp string, signature string) bool {

	// 获取 path params
	params := map[string]interface{}{}
	for k, v := range mux.Vars(r) {
		params[k] = v
	}

	// 获取 query params
	for k, v := range r.URL.Query() {
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
	hexHash := hash(timestamp + appKey)

	for k, v := range paramsBody {
		params[k] = v
	}

	// sort params
	sortParams := sortMapParams(params)

	if sortParams != nil {
		sortParamsBytes, _ := json.Marshal(sortParams)
		hexHash = hash(string(sortParamsBytes) + timestamp + appKey)
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
