package controller

import (
	"context"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service/ddc"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service/native"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/chain"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/handlers"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/kit"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
)

const namespace = "nftp"

type (
	//BaseController define a base controller for all http Controller
	BaseController struct {
		kit.Controller
	}
)

// GetAllControllers return all the controllers of the app server
func GetAllControllers() []kit.IController {
	bc := BaseController{
		Controller: kit.NewController(),
	}

	baseSvc := service.NewBase(chain.GetSdkClient(), chain.GetGas(), chain.GetDenom(), chain.GetAmount())
	controllers := []kit.IController{
		NewDemoController(bc, handlers.NewDemo()),
		NewAccountsController(bc, handlers.NewAccount(native.NewNFTAccount(baseSvc), ddc.NewDDCAccount(baseSvc))),
		NewNftClassController(bc, handlers.NewNFTClass(native.NewNFTClass(baseSvc))),
		NewNftController(bc, handlers.NewNft(native.NewNFT(baseSvc))),
		NewNftTransferController(bc, handlers.NewNftTransfer(native.NewNftTransfer(baseSvc))),
		NewTxController(bc, handlers.NewTx(service.NewTx())),
	}

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
