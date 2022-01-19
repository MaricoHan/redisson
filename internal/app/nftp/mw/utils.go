package mw

import (
	"encoding/json"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

func writeBadRequestResp(w http.ResponseWriter, err types.IError) {
	w.WriteHeader(http.StatusBadRequest)
	response := kit.Response{
		ErrorResp: &kit.ErrorResp{
			Code:    err.Code(),
			Message: err.Error(),
		},
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}

func writeForbiddenResp(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	response := kit.Response{
		ErrorResp: &kit.ErrorResp{
			Code:    types.ErrParams.Code(),
			Message: types.ErrParams.Error(),
		},
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}
