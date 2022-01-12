package controller

import (
	"context"

	"gitlab.bianjie.ai/irita-nftp/nftp-open-api/internal/pkg/kit"
	"gitlab.bianjie.ai/irita-nftp/nftp-open-api/internal/pkg/log"
)

const namespace = "nftp"

type (
	//BaseController define a base controller for all http Controller
	BaseController struct {
		kit.Controller
	}
)

//return all the controllers of the app server
func GetAllControllers() []kit.IController {
	_ = BaseController{
		Controller: kit.NewController(),
	}

	controllers := []kit.IController{}

	return controllers
}

// makeHandler create a http hander for request
func (bc BaseController) makeHandler(h kit.Handler, request interface{}) *kit.Server {
	return bc.MakeHandler(
		bc.wrapHandler(h),
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
			log.Error("Execute handler logic failed", "error", err.Error())
			return nil, err
		}
		return resp, nil
	}
}
