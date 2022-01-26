package dto

import (
	"github.com/volatiletech/sqlboiler/v4/types"
)

type CreateAccountP struct {
	Count int64  `json:"count"`
	AppID uint64 `json:"app_id"`
}

type AccountsP struct {
	PageP
	Account   string `json:"account"`
	AppID     uint64 `json:"app_id"`
	Module    string `json:"module"`
	Operation string `json:"operation"`
}

type AccountsRes struct {
	PageRes
	Accounts []*Account `json:"accounts"`
}

type Account struct {
	Account string `json:"account"`
	Gas     uint64 `json:"gas"`
}

type AccountOperationRecordRes struct {
	PageRes
	OperationRecords []*AccountOperationRecords `json:"operation_records"`
}

type AccountOperationRecords struct {
	TxHash    string     `json:"tx_hash"`
	Module    string     `json:"module"`
	Operation string     `json:"operation"`
	Signer    string     `json:"signer"`
	Timestamp string     `json:"timestamp"`
	Message   types.JSON `json:"message"`
}
