package controller

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
)

// BaseController define a base controller for all http Controller
type BaseController struct {
	Controller kit.Controller
}

func GetAllControllers(logger *log.Logger) []kit.IController {
	baseController := BaseController{
		Controller: kit.NewController(),
	}
	controllers := []kit.IController{
		NewAccountsController(baseController, handlers.NewAccount(services.NewAccount(logger))),
		NewMsgsController(baseController, handlers.NewMsgs(services.NewMsgs(logger))),
		NewTxController(baseController, handlers.NewTx(services.NewTx(logger))),
		NewNftClassController(baseController, handlers.NewNFTClass(services.NewNFTClass(logger))),
		NewNftController(baseController, handlers.NewNft(services.NewNFT(logger))),
		NewNftTransferController(baseController, handlers.NewNFTTransfer(services.NewNFTTransfer(logger))),
		NewAuthController(baseController, handlers.NewAuth(services.NewAuth(logger))),
		NewUserController(baseController, handlers.NewUser(services.NewUser(logger))),
		NewNsController(baseController, handlers.NewNs(services.NewNs(logger))),
		NewRecordController(baseController, handlers.NewRecord(services.NewRecord(logger))),
		NewEmptionController(baseController, handlers.NewBusiness(services.NewBusiness(logger))),
		NewContractController(baseController, handlers.NewContract(services.NewContract(logger))),
	}

	return controllers
}

// makeHandler create a http hander for request
func (bc BaseController) makeHandler(handler kit.Handler, request interface{}) *kit.Server {
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
