package app

import (
	"encoding/json"
	"fmt"

	"net/http"
	"runtime/debug"

	"gitlab.bianjie.ai/irita-nftp/nftp-open-api/internal/pkg/kit"
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
			w.WriteHeader(500)
			bz, _ := json.Marshal(kit.Response{
				Code:    kit.Failed,
				Message: "system panic",
			})
			_, _ = w.Write(bz)
		}
	}()
	h.next.ServeHTTP(w, r)
}
