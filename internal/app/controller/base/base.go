package base

import (
	"context"

	log "github.com/sirupsen/logrus"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

// BaseController define a base controller for all http Controller
type BaseController struct {
	Controller kit.Controller
}

// MakeHandler create a http hander for request
func (bc BaseController) MakeHandler(handler kit.Handler, request interface{}) *kit.Server {
	return bc.Controller.MakeHandler(
		bc.wrapHandler(handler),
		request,
		[]kit.RequestFunc{},
		nil,
		[]kit.ServerResponseFunc{},
	)
}

func (bc BaseController) wrapHandler(h kit.Handler) kit.Handler {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Debug("Execute handler logic ", "method", "wrapHandler", "params", request)
		resp, err := h(ctx, request)
		if err != nil {
			// log.Error("Execute handler logic failed", "error", err.Error())
			return nil, err
		}
		return resp, nil
	}
}
