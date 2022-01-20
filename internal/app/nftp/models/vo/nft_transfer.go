package vo

import "go/types"

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
	Recipients types.Object[] `json:"recipients{index,recipient}" validate:"required"`
}
