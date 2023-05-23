package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/notice"
)

type NoticeController struct {
	base.BaseController
	handler handlers.INotice
}

func NewNoticeController(bc base.BaseController, handler handlers.INotice) kit.IController {
	return NoticeController{bc, handler}
}

func (c NoticeController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/notice/nfts",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferNFTS, &vo.TransferNFTS{}),
		},
		kit.Endpoint{
			URI:     "/notice/classes",
			Method:  http.MethodPost,
			Handler: c.MakeHandler(c.handler.TransferClasses, &vo.TransferClasses{}),
		},
	)
	return ends
}
