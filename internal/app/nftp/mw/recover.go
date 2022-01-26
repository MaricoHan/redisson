package mw

import (
	"encoding/json"
	"fmt"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"

	"net/http"
	"runtime/debug"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
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
			bz, _ := json.Marshal(kit.Response{
				ErrorResp: &kit.ErrorResp{
					CodeSpace: types.ErrInternal.CodeSpace(),
					Code:      types.ErrInternal.Code(),
					Message:   types.ErrInternal.Error(),
				},
			})
			_, _ = w.Write(bz)
		}
	}()
	h.next.ServeHTTP(w, r)
}
