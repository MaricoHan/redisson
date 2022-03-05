package vo

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

type TransferNftClassByIDRequest struct {
	Base
	Recipient string                 `json:"recipient" validate:"required"`
	Tag       map[string]interface{} `json:"tag"`
}

type TransferNftByNftIdRequest struct {
	Base
	Recipient string                 `json:"recipient" validate:"required"`
	Tag       map[string]interface{} `json:"tag"`
}

type TransferNftByBatchRequest struct {
	Base
	Recipients []*dto.Recipient `json:"recipients" validate:"required"`
}
