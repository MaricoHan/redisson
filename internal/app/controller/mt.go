package controller

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/mt"
	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"
	"net/http"
)

type MTController struct {
	BaseController
	handler handlers.IMT
}

func NewMTController(bc BaseController, handler handlers.IMT) kit.IController {
	return MTController{bc, handler}
}

func (m MTController) GetEndpoints() []kit.Endpoint {
	return []kit.Endpoint{
		{
			URI:     "/mt/mts-issue/{class_id}",
			Method:  http.MethodPost,
			Handler: m.makeHandler(m.handler.Issue, &vo.IssueRequest{}),
		},
		{
			URI:     "/mt/mts-mint/{class_id}/{mt_id}",
			Method:  http.MethodPost,
			Handler: m.makeHandler(m.handler.Mint, &vo.IssueRequest{}),
		},
	}
}
