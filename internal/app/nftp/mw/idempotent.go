package mw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/redis"
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
		writeBadRequestResp(w, types.ErrParams)
		return
	}
	appID := r.Header.Get("X-App-Id")
	key := fmt.Sprintf("%s:%s", appID, req.OperationID)
	ok, err := redis.Has(key)
	if err != nil {
		writeInternalResp(w)
		return
	}
	if ok {
		writeBadRequestResp(w, types.ErrIdempotent)
		return
	}

	if err := redis.Set(key, "1", time.Second*60); err != nil {
		writeBadRequestResp(w, types.ErrRedisConn)
		return
	}

	h.next.ServeHTTP(w, r)
}
