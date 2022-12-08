package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"gitlab.bianjie.ai/avata/open-api/internal/app/models/entity"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
	"gitlab.bianjie.ai/avata/open-api/utils"
)

// RouterMiddleware recover the panic error
func RouterMiddleware(h http.Handler) http.Handler {
	return routerHandler{h}
}

type routerHandler struct {
	next http.Handler
}

func (router routerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := initialize.Log.WithFields(map[string]interface{}{
		"function": "ServeHTTP",
		"method":   r.Method,
		"url":      r.URL.Path,
	})
	authData, err := utils.HeaderAuthData(&r.Header)
	if err != nil {
		log.WithError(err).Error("get auth data")
		writeInternalResp(w)
		return
	}
	// DDC 不支持 NFT-批量、orders-批量、MT、版权服务
	if fmt.Sprintf("%s-%s", authData.Code, authData.Module) == constant.WenchangDDC {
		if strings.Contains(r.RequestURI, "/mt/") || strings.Contains(r.RequestURI, "/nft/batch/") || strings.Contains(r.RequestURI, "/orders/batch") || strings.Contains(r.RequestURI, "/rights/") {
			writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
			return
		}
	} else { // native
		// 非托管模式
		if authData.AccessMode == entity.UNMANAGED {
			if strings.Contains(r.RequestURI, "/rights/") {
				writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
				return
			}
			if fmt.Sprintf("%s-%s", authData.Code, authData.Module) == constant.IritaOPBNative {
				// 文昌链-天舟除 orders 都不支持
				if !strings.Contains(r.RequestURI, "/orders") && !strings.Contains(r.RequestURI, "/auth") {
					writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
					return
				}
			} else if !strings.Contains(r.RequestURI, "/auth") {
				// 文昌链-天和都不支持
				writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
				return
			}

		} else {
			// 托管不支持 orders
			if strings.Contains(r.RequestURI, "/orders") {
				writeNotFoundRequestResp(w, constant.ErrUnmanagedUnSupported)
				return
			}
		}
	}
	router.next.ServeHTTP(w, r)
}
