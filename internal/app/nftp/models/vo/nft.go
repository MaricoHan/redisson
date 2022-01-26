package vo

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

type EditNftByIndexRequest struct {
	Name string `json:"name" validate:"required"`
	Uri  string `json:"uri"`
	Data string `json:"data"`
}

type EditNftByBatchRequest struct {
	EditNftsR []*dto.EditNft `json:"nfts"`
}
