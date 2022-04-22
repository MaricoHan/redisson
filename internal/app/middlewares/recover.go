package middlewares

import (
	"encoding/json"
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"net/http"
	"runtime/debug"
)

// RecoverMiddleware recover the panic error
func RecoverMiddleware(h http.Handler) http.Handler {
	return panicHandler{h}
}

type panicHandler struct {
	next http.Handler
}

func (h panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e, string(debug.Stack()))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			bz, _ := json.Marshal(constant.Response{
				ErrorResp: &constant.ErrorResp{
					CodeSpace: constant.ErrInternal.CodeSpace(),
					Code:      constant.ErrInternal.Code(),
					Message:   constant.ErrInternal.Error(),
				},
			})
			_, _ = w.Write(bz)
		}
	}()
	h.next.ServeHTTP(w, r)
}
