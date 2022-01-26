package mw

import (
	"encoding/json"
	"net/http"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

func writeBadRequestResp(w http.ResponseWriter, err types.IError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	response := kit.Response{
		ErrorResp: &kit.ErrorResp{
			CodeSpace: err.CodeSpace(),
			Code:      err.Code(),
			Message:   err.Error(),
		},
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}

func writeForbiddenResp(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	response := kit.Response{
		ErrorResp: &kit.ErrorResp{
			CodeSpace: types.ErrParams.CodeSpace(),
			Code:      types.ErrParams.Code(),
			Message:   types.ErrParams.Error(),
		},
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}

func writeInternalResp(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	response := kit.Response{
		ErrorResp: &kit.ErrorResp{
			CodeSpace: types.ErrInternal.CodeSpace(),
			Code:      types.ErrInternal.Code(),
			Message:   types.ErrInternal.Error(),
		},
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}
