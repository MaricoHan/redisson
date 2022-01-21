package vo

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

type EditNftByIndexRequest struct {
	Base
	Name string `json:"name" validate:"required"`
	Uri  string `json:"uri"`
	Data string `json:"data"`
}

type EditNftByBatchRequest struct {
	Base
	EditNftsR []*dto.EditNft `json:"edit_nfts"`
}
