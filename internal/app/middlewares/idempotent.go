package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"
	"io/ioutil"
	"net/http"
	"time"

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
	if r.Method != http.MethodPost {
		h.next.ServeHTTP(w, r)
		return
	}
	// 把request的内容读取出来
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	// 把刚刚读出来的再写进去
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	req := &vo.Base{}
	err := json.Unmarshal(bodyBytes, req)
	if err != nil {
		log.Error("server http", "params error:", err)
		writeBadRequestResp(w, constant.ErrParams)
		return
	}
	if len(req.OperationID) == 0 {
		writeBadRequestResp(w, constant.NewAppError(constant.RootCodeSpace, errors2.StrToCode[errors2.RequestsFailed], "operation_id is a required field"))
		return
	}
	if len(req.OperationID) >= 65 {
		writeBadRequestResp(w, constant.NewAppError(constant.RootCodeSpace, errors2.StrToCode[errors2.RequestsFailed], "operation_id does not comply with the rules"))
		return
	}

	appID := r.Header.Get("X-App-Id")
	key := fmt.Sprintf("%s:%s", appID, req.OperationID)
	ok, err := initialize.RedisClient.Has(key)
	if err != nil {
		log.Error("redis error", "redis get error:", err)
		writeInternalResp(w)
		return
	}

	if ok {
		writeBadRequestResp(w, constant.ErrIdempotent)
		return
	}
	if err := initialize.RedisClient.Set(key, "1", time.Second*60); err != nil {
		log.Error("redis error", "redis set error:", err)
		writeBadRequestResp(w, constant.ErrInternal)
		return
	}
	w.Header().Set("X-Operation-ID", req.OperationID)
	h.next.ServeHTTP(w, r)
}
