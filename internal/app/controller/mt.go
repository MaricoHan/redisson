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
			URI:     "/mt/mt-issues/{class_id}",
			Method:  http.MethodPost,
			Handler: m.makeHandler(m.handler.Issue, &vo.IssueRequest{}),
		},
		{
			URI:     "/mt/mt-mints/{class_id}/{mt_id}",
			Method:  http.MethodPost,
			Handler: m.makeHandler(m.handler.Mint, &vo.MintRequest{}),
		},
		{
			URI:     "/mt/mts",
			Method:  http.MethodGet,
			Handler: m.makeHandler(m.handler.List, nil),
		},
		{
			URI:     "/mt/mts/{class_id}/{mt_id}",
			Method:  http.MethodGet,
			Handler: m.makeHandler(m.handler.Show, nil),
		},
		{
			URI:     "/mt/mts/{class_id}/{account}/balances",
			Method:  http.MethodGet,
			Handler: m.makeHandler(m.handler.Balances, nil),
		},
		{
			URI:     "/mt/mts/{class_id}/{owner}/{mt_id}",
			Method:  http.MethodPatch,
			Handler: m.makeHandler(m.handler.Edit, &vo.EditRequest{}),
		},
		//{
		//	URI:     "/mt/mts/{owner}",
		//	Method:  http.MethodDelete,
		//	Handler: m.makeHandler(m.handler.Burn, &vo.BatchBurnRequest{}),
		//},
		{
			URI:     "/mt/mts/{class_id}/{owner}/{mt_id}",
			Method:  http.MethodDelete,
			Handler: m.makeHandler(m.handler.Burn, &vo.BurnRequest{}),
		},
		{
			URI:     "/mt/mt-transfers/{class_id}/{owner}/{mt_id}",
			Method:  http.MethodPost,
			Handler: m.makeHandler(m.handler.Transfer, &vo.TransferRequest{}),
		},
	}
}