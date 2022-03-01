package vo

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

type CreateNftsRequest struct {
	Base
	Name      string `json:"name" validate:"required"`
	Uri       string `json:"uri"`
	UriHash   string `json:"uri_hash"`
	Data      string `json:"data"`
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
	Tag       string `json:"tag"`
}

type EditNftByIndexRequest struct {
	Name string `json:"name" validate:"required"`
	Uri  string `json:"uri"`
	Data string `json:"data"`
	Tag  string `json:"tag"`
}

type EditNftByBatchRequest []*dto.EditNft
