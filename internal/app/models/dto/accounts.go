package dto

import (
	"github.com/volatiletech/sqlboiler/types"
)

type CreateAccount struct {
	Count      int64  `json:"count"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
}

type AccountsInfo struct {
	Page
	Account         string `json:"account"`
	ProjectID       uint64 `json:"project_id"`
	ChainID         uint64 `json:"chain_id"`
	PlatFormID      uint64 `json:"plat_form_id"`
	Module          string `json:"module"`
	Operation       string `json:"operation"`
	OperationModule string `json:"operation_module"`
	Code            string `json:"code"`
}

type AccountRes struct {
	Accounts []string `json:"accounts"`
}

type AccountsRes struct {
	PageRes
	Accounts []*Account `json:"accounts"`
}

type Account struct {
	Account string `json:"account"`
	Gas     uint64 `json:"gas"`
	BizFee  uint64 `json:"biz_fee"` // 余额业务
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
