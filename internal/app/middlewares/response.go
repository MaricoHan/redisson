package middlewares

import (
	"encoding/json"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"net/http"
)

func writeBadRequestResp(w http.ResponseWriter, err constant.IError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	response := constant.Response{
		ErrorResp: &constant.ErrorResp{
			CodeSpace: err.CodeSpace(),
			Code:      err.Code(),
			Message:   err.Error(),
		},
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}

func writeNotFoundRequestResp(w http.ResponseWriter, err constant.IError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	response := constant.Response{
		ErrorResp: &constant.ErrorResp{
			CodeSpace: err.CodeSpace(),
			Code:      err.Code(),
			Message:   err.Error(),
		},
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}

func writeForbiddenResp(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	response := constant.Response{
		ErrorResp: &constant.ErrorResp{
			CodeSpace: constant.ErrAuthenticate.CodeSpace(),
			Code:      constant.ErrAuthenticate.Code(),
			Message:   constant.ErrAuthenticate.Error(),
		},
	}
	if errMsg != "" {
		response.ErrorResp.Message = errMsg
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}

func writeInternalResp(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	response := constant.Response{
		ErrorResp: &constant.ErrorResp{
			CodeSpace: constant.ErrInternal.CodeSpace(),
			Code:      constant.ErrInternal.Code(),
			Message:   constant.ErrInternal.Error(),
		},
	}
	bz, _ := json.Marshal(response)
	_, _ = w.Write(bz)
	return
}
