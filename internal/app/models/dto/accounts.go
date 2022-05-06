package dto

import (
	"github.com/volatiletech/sqlboiler/types"
)

// BatchCreateAccount 批量创建链账户
type BatchCreateAccount struct {
	Count       int64  `json:"count"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	OperationId string `json:"operation_id"`
}

// CreateAccount 创建链账户
type CreateAccount struct {
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	OperationId string `json:"operation_id"`
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
	TxHash          string `json:"tx_hash"`
	OperationId     string `json:"operation_id"`
}

type BatchAccountRes struct {
	Accounts    []string `json:"accounts"`
	OperationId string   `json:"operation_id"`
}

type AccountRes struct {
	Account     string `json:"account"`
	Name        string `json:"name"`
	OperationId string `json:"operation_id"`
}

type AccountsRes struct {
	PageRes
	Accounts []*Account `json:"accounts"`
}

type Account struct {
	Account     string `json:"account"`
	Name        string `json:"name"`
	OperationId string `json:"operation_id"`
	Gas         uint64 `json:"gas"`
	BizFee      uint64 `json:"biz_fee"` // 余额业务
}

type AccountOperationRecordRes struct {
	PageRes
	OperationRecords []*AccountOperationRecords `json:"operation_records"`
}

type AccountOperationRecords struct {
	TxHash      string     `json:"tx_hash"`
	Module      string     `json:"module"`
	Operation   string     `json:"operation"`
	Signer      string     `json:"signer"`
	Timestamp   string     `json:"timestamp"`
	GasFee      uint64     `json:"gas_fee"`
	BusinessFee uint64     `json:"business_fee"`
	Message     types.JSON `json:"message"`
}
