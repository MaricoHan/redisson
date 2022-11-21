package controller

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/notice"
)

type NoticeController struct {
	BaseController
	handler handlers.INotice
}

func NewNoticeController(bc BaseController, handler handlers.INotice) kit.IController {
	return NoticeController{bc, handler}
}

func (c NoticeController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends,
		kit.Endpoint{
			URI:     "/notice/nfts",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferNFTS, &notice.TransferNFTS{}),
		},
		kit.Endpoint{
			URI:     "/notice/classes",
			Method:  http.MethodPost,
			Handler: c.makeHandler(c.handler.TransferClasses, &notice.TransferClasses{}),
		},
	)
	return ends
}
