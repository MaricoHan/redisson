package mw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/redis"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
	"io/ioutil"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
)

func IdempotentMiddleware(h http.Handler) http.Handler {
	return idempotentMiddlewareHandler{h}
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
		writeBadRequestResp(w, types.ErrParams)
		return
	}
	if len(req.OperationID) >= 65 || len(req.OperationID) == 0 {
		writeBadRequestResp(w, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, "operation_id does not comply with the rules"))
		return
	}

	appID := r.Header.Get("X-App-Id")
	key := fmt.Sprintf("%s:%s", appID, req.OperationID)
	ok, err := redis.Has(key)
	if err != nil {
		log.Error("redis error", "redis error:", err)
		writeInternalResp(w)
		return
	}

	if ok {
		writeBadRequestResp(w, types.ErrIdempotent)
		return
	}
	w.Header().Set("X-Operation-ID", req.OperationID)
	h.next.ServeHTTP(w, r)
}
