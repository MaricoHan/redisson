package mw

import (
	"net/http"
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
	r.Header.Set("X-App-ID", "sheldon test")
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
	//// appID := r.Header.Get("X-App-ID")
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
	//// 增加
	//w.Header().Add("X-App-ID", "")

	h.next.ServeHTTP(w, r)
}

func (h authHandler) Signature() {

}
