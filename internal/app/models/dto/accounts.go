package dto

import (
	"github.com/volatiletech/sqlboiler/types"
)

// BatchCreateAccount 批量创建链账户
type BatchCreateAccount struct {
	Count       uint32 `json:"count"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	OperationId string `json:"operation_id"`
	AccessMode  int    `json:"access_mode"`
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
	AccessMode  int    `json:"access_mode"`
}

type AccountsInfo struct {
	Page
	Account         string `json:"account"`
	ProjectID       uint64 `json:"project_id"`
	ChainID         uint64 `json:"chain_id"`
	PlatFormID      uint64 `json:"plat_form_id"`
	Module          string `json:"module"`
	Operation       uint32 `json:"operation"`
	OperationModule uint32 `json:"operation_module"`
	Code            string `json:"code"`
	TxHash          string `json:"tx_hash"`
	OperationId     string `json:"operation_id"`
	Name            string `json:"name"`
	AccessMode      int    `json:"access_mode"`
}

type BatchAccountRes struct {
	Accounts []string `json:"accounts"`
}

type AccountRes struct {
	Account string `json:"account"`
}

type AccountsRes struct {
	PageRes
	Accounts []*Account `json:"accounts"`
}

type Account struct {
	Account     string `json:"account"`
	Name        string `json:"name"`
	OperationId string `json:"operation_id"`
}

type AccountOperationRecordRes struct {
	PageRes
	OperationRecords []*AccountOperationRecords `json:"operation_records"`
}

type AccountOperationRecords struct {
	TxHash    string      `json:"tx_hash"`
	Module    uint32      `json:"module"`
	Operation uint32      `json:"operation"`
	Signer    string      `json:"signer"`
	Timestamp string      `json:"timestamp"`
	NftMsg    *types.JSON `json:"nft_msg"`
}
