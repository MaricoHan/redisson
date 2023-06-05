package controller

import (
	log "github.com/sirupsen/logrus"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	evm_controller "gitlab.bianjie.ai/avata/open-api/internal/app/controller/evm"
	l2_controller "gitlab.bianjie.ai/avata/open-api/internal/app/controller/l2"
	native_controller "gitlab.bianjie.ai/avata/open-api/internal/app/controller/native"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/evm"
	l2_handlers "gitlab.bianjie.ai/avata/open-api/internal/app/handlers/l2"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services"
	evm2 "gitlab.bianjie.ai/avata/open-api/internal/app/services/evm"
	l2_services "gitlab.bianjie.ai/avata/open-api/internal/app/services/l2"
	native2 "gitlab.bianjie.ai/avata/open-api/internal/app/services/native"
)

func GetAllControllers(logger *log.Logger) []kit.IController {
	baseController := base.BaseController{
		Controller: kit.NewController(),
	}
	controllers := []kit.IController{
		// 公共接口
		NewAccountsController(baseController, handlers.NewAccount(services.NewAccount(logger))),
		NewAuthController(baseController, handlers.NewAuth(services.NewAuth(logger))),
		NewUserController(baseController, handlers.NewUser(services.NewUser(logger))),
		NewEmptionController(baseController, handlers.NewBusiness(services.NewBusiness(logger))),
		NewNoticeController(baseController, handlers.NewNotice(services.NewNotice(logger))),
		NewRecordController(baseController, handlers.NewRecord(services.NewRecord(logger))),
		// layer接口
		l2_controller.NewNftClassController(baseController, l2_handlers.NewNFTClass(l2_services.NewNFTClass(logger))),
		l2_controller.NewNftController(baseController, l2_handlers.NewNft(l2_services.NewNFT(logger))),
		l2_controller.NewDictController(baseController, l2_handlers.NewDict(l2_services.NewDict(logger))),
		// EVM接口
		evm_controller.NewNsController(baseController, evm.NewNs(evm2.NewNs(logger))),
		evm_controller.NewMsgsController(baseController, evm.NewMsgs(evm2.NewMsgs(logger))),
		evm_controller.NewNftClassController(baseController, evm.NewNFTClass(evm2.NewNFTClass(logger))),
		evm_controller.NewNftController(baseController, evm.NewNft(evm2.NewNFT(logger))),
		evm_controller.NewContractController(baseController, evm.NewContract(evm2.NewContract(logger))),
		evm_controller.NewNftTransferController(baseController, evm.NewNFTTransfer(evm2.NewNFTTransfer(logger))),
		evm_controller.NewDictController(baseController, evm.NewDict(evm2.NewDict(logger))),
		evm_controller.NewTxController(baseController, evm.NewTx(evm2.NewTx(logger))),

		// native接口
		native_controller.NewMTClassController(baseController, native.NewMTClass(native2.NewMTClass(logger))),
		native_controller.NewMTController(baseController, native.NewMT(native2.NewMT(logger))),
		native_controller.NewMsgsController(baseController, native.NewMsgs(native2.NewMsgs(logger))),
		native_controller.NewNftController(baseController, native.NewNft(native2.NewNft(logger))),
		native_controller.NewNftClassController(baseController, native.NewNFTClass(native2.NewNFTClass(logger))),
		native_controller.NewNFTTransferController(baseController, native.NewNFTTransfer(native2.NewNFTTransfer(logger))),
		native_controller.NewDictController(baseController, native.NewDict(native2.NewDict(logger))),
		native_controller.NewTxController(baseController, native.NewTx(native2.NewTx(logger))),
	}

	return controllers
}
