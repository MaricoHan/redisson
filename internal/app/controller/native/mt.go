package native

import (
	"net/http"

	kit "gitlab.bianjie.ai/avata/open-api/pkg/gokit"

	"gitlab.bianjie.ai/avata/open-api/internal/app/controller/base"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers/native"
	vo "gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/mt"
)

type MTController struct {
	base.BaseController
	handler native.IMT
}

func NewMTController(bc base.BaseController, handler native.IMT) kit.IController {
	return MTController{bc, handler}
}

func (m MTController) GetEndpoints() []kit.Endpoint {
	return []kit.Endpoint{
		{
			URI:     "/native/mt/mt-issues/{class_id}",
			Method:  http.MethodPost,
			Handler: m.MakeHandler(m.handler.Issue, &vo.IssueRequest{}),
		},
		{
			URI:     "/native/mt/mt-mints/{class_id}/{mt_id}",
			Method:  http.MethodPost,
			Handler: m.MakeHandler(m.handler.Mint, &vo.MintRequest{}),
		},
		{
			URI:     "/native/mt/mts",
			Method:  http.MethodGet,
			Handler: m.MakeHandler(m.handler.List, nil),
		},
		{
			URI:     "/native/mt/mts/{class_id}/{mt_id}",
			Method:  http.MethodGet,
			Handler: m.MakeHandler(m.handler.Show, nil),
		},
		{
			URI:     "/native/mt/mts/{class_id}/{account}/balances",
			Method:  http.MethodGet,
			Handler: m.MakeHandler(m.handler.Balances, nil),
		},
		{
			URI:     "/native/mt/mts/{class_id}/{owner}/{mt_id}",
			Method:  http.MethodPatch,
			Handler: m.MakeHandler(m.handler.Edit, &vo.EditRequest{}),
		},
		//{
		//	URI:     "/mt/mts/{owner}",
		//	Method:  http.MethodDelete,
		//	Handler: m.MakeHandler(m.handler.Burn, &vo.BatchBurnRequest{}),
		//},
		{
			URI:     "/native/mt/mts/{class_id}/{owner}/{mt_id}",
			Method:  http.MethodDelete,
			Handler: m.MakeHandler(m.handler.Burn, &vo.BurnRequest{}),
		},
		{
			URI:     "/native/mt/mt-transfers/{class_id}/{owner}/{mt_id}",
			Method:  http.MethodPost,
			Handler: m.MakeHandler(m.handler.Transfer, &vo.TransferRequest{}),
		},
	}
}
