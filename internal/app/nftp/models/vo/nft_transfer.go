package vo

import "gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"

type TransferNftClassByID struct {
	Base
	Recipient string `json:"recipient" validate:"required"`
}

type TransferNftByIndex struct {
	Base
	Recipient string `json:"recipient" validate:"required"`
}

type TransferNftByBatch struct {
	Base
	Recipients []*dto.Recipient `json:"recipients" validate:"required"`
}
