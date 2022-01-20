package mw

import (
	"context"
	"fmt"
	"net/http"

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

	appKey := r.Header.Get("X-Api-Key")
	appKeyResult, err := models.TAppKeys(
		qm.Select(models.TAppKeyColumns.AppID),
		models.TAppKeyWhere.APIKey.EQ(appKey),
	).OneG(context.Background())
	if err != nil {
		writeForbiddenResp(w)
		return
	}

	r.Header.Set("X-App-Id", fmt.Sprintf("%d", appKeyResult.AppID))
	//// 1. 获取 header 中的时间戳
	//reqTimestampStr := r.Header.Get("X-Timestamp")
	//reqTimestampInt, err := strconv.ParseInt(reqTimestampStr, 10, 64)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//// 1.1 判断时间误差
	//curTimestamp := time.Now().Unix()
	//if curTimestamp-reqTimestampInt > timeInterval || curTimestamp < reqTimestampInt {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//// 如果时间误差超过
	////2. 获取 header 中的 app_id
	//// 	从数据中查询用户的信息
	//// params + reqTimestampStr + appKey
	//// appID := r.Header.Get("X-App-Id")
	//switch r.Method {
	//case http.MethodGet:
	//case http.MethodPost:
	//case http.MethodPatch:
	//case http.MethodDelete:
	//	//defer r.Body.Close()
	//	//body, err := ioutil.ReadAll(r.Body)
	//	//if err != nil {
	//	//	w.WriteHeader(http.StatusBadRequest)
	//	//	return
	//	//}
	//default:
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//}

	h.next.ServeHTTP(w, r)
}

func (h authHandler) Signature() {

}
