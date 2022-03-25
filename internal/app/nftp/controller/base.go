package controller

import (
	"context"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/handlers"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service/wenchangchain-ddc"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service/wenchangchain-native"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/chain"
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
	baseSecs := make(map[string]*service.Base)
	for k, v := range chain.GetSdkClients() {
		baseSvc := service.NewBase(v.Client, v.Gas, v.Denom, v.Amount)
		baseSecs[k] = baseSvc
	}
	controllers := []kit.IController{
		NewDemoController(bc, handlers.NewDemo()),
		NewAccountsController(bc, handlers.NewAccount(wenchangchain_native.NewNFTAccount(baseSecs), wenchangchain_ddc.NewDDCAccount(baseSecs))),
		NewNftClassController(bc, handlers.NewNFTClass(wenchangchain_native.NewNFTClass(baseSecs), wenchangchain_ddc.NewDDCClass(baseSecs))),
		NewNftController(bc, handlers.NewNft(wenchangchain_native.NewNFT(baseSecs), wenchangchain_ddc.NewDDC(baseSecs))),
		NewNftTransferController(bc, handlers.NewNftTransfer(wenchangchain_native.NewNftTransfer(baseSecs), wenchangchain_ddc.NewDDCTransfer(baseSecs))),
		NewTxController(bc, handlers.NewTx(wenchangchain_native.NewTx(), wenchangchain_ddc.NewTx())),
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
