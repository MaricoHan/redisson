package controller

import (
	log "github.com/sirupsen/logrus"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	l2_controller "gitlab.bianjie.ai/avata/open-api/internal/app/controller/l2"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	l2_handlers "gitlab.bianjie.ai/avata/open-api/internal/app/handlers/l2"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	l2_services "gitlab.bianjie.ai/avata/open-api/internal/app/services/l2"
)

func GetAllControllers(logger *log.Logger) []kit.IController {
	baseController := base.BaseController{
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
		NewContractController(baseController, handlers.NewContract(services.NewContract(logger))),
		l2_controller.NewNftClassController(baseController, l2_handlers.NewNFTClass(l2_services.NewNFTClass(logger))),
	}

	return controllers
}
