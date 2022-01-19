package mw

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeBadRequestResp(w, types.ErrParams)
		return
	}
	req := &vo.Base{}
	err = json.Unmarshal(body, req)
	if err != nil {
		writeBadRequestResp(w, types.ErrParams)
		return
	}

	if redis.Has(req.OperationID) {
		writeBadRequestResp(w, types.ErrIdempotent)
		return
	}

	h.next.ServeHTTP(w, r)
}
