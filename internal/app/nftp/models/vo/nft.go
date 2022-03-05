package vo

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

type CreateNftsRequest struct {
	Base
	Name    string `json:"name" validate:"required"`
	Uri     string `json:"uri"`
	UriHash string `json:"uri_hash"`
	Data    string `json:"data"`
	//关闭批量发行
	//Amount    int    `json:"amount"`
	Recipient string                 `json:"recipient"`
	Tag       map[string]interface{} `json:"tag"`
}

type EditNftByIndexRequest struct {
	Name string                 `json:"name" validate:"required"`
	Uri  string                 `json:"uri"`
	Data string                 `json:"data"`
	Tag  map[string]interface{} `json:"tag"`
}

type EditNftByBatchRequest []*dto.EditNft

type DeleteNftByNftIdRequest struct {
	Tag map[string]interface{} `json:"tag"`
}
