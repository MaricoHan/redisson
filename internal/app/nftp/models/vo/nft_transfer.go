package vo

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

type TransferNftClassByIDRequest struct {
	Base
	Recipient string `json:"recipient" validate:"required"`
}

type TransferNftByIndexRequest struct {
	Base
	Recipient string `json:"recipient" validate:"required"`
}

type TransferNftByBatchRequest struct {
	Base
	Recipients []*dto.Recipient `json:"recipients" validate:"required"`
}
