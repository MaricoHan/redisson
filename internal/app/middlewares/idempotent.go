package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"

	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo"
)

func IdempotentMiddleware(handler http.Handler) http.Handler {
	return idempotentMiddlewareHandler{next: handler}
}

type idempotentMiddlewareHandler struct {
	next http.Handler
}

func (h idempotentMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 把request的内容读取出来
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	if bodyBytes == nil || len(bodyBytes) < 1 {
		h.next.ServeHTTP(w, r)
		return
	}
	// if r.Method != http.MethodPost {
	//	h.next.ServeHTTP(w, r)
	//	return
	// }

	// 把刚刚读出来的再写进去
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	req := &vo.Base{}
	err := json.Unmarshal(bodyBytes, req)
	if err != nil {
		log.Error("server http params error:", err)
		writeBadRequestResp(w, constant.ErrParams)
		return
	}

	if len(req.OperationID) < 1 {
		order := &vo.Order{}
		err = json.Unmarshal(bodyBytes, order)
		if err != nil {
			log.Error("server http params error:", err)
			writeBadRequestResp(w, constant.ErrParams)
			return
		}
		if len(order.OperationID) < 1 {
			// 部分接口 operation_id 不是必填
			// 不存在 operation_id 请求，具体验证规则由目标微服务处理
			h.next.ServeHTTP(w, r)
			return
		}
		req.OperationID = order.OperationID
	}
	// if len(req.OperationID) >= 65 {
	//	writeBadRequestResp(w, constant.NewAppError(constant.RootCodeSpace, errors2.StrToCode[errors2.DuplicateRequest], "operation_id does not comply with the rules"))
	//	return
	// }

	appID := r.Header.Get("X-App-Id")
	key := fmt.Sprintf("%s:%s", appID, req.OperationID)
	incr := initialize.RedisClient.Incr(key)
	if incr > 1 {
		writeBadRequestResp(w, constant.ErrDuplicate)
		return
	}
	if err := initialize.RedisClient.Expire(key, time.Second*time.Duration(constant.GrpcTimeout)); err != nil {
		// 自动过期时间设置大于open-api超时时间
		log.Error("redis error", "redis set error:", err)
		writeBadRequestResp(w, constant.ErrInternal)
		return
	}
	w.Header().Set("X-Operation-ID", req.OperationID)
	r.Header.Set("X-App-Operation-Key", fmt.Sprintf("%s", key))
	h.next.ServeHTTP(w, r)
}
