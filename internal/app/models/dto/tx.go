package dto

import "github.com/volatiletech/sqlboiler/types"

type TxResultByTxHash struct {
	TaskId     string `json:"task_id"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
}

type TxResultByTxHashRes struct {
	Module      string                 `json:"module"`
	Type        string                 `json:"type"`
	TxHash      string                 `json:"tx_hash"`
	Status      int32                  `json:"status"`
	Nft         *types.JSON            `json:"nft"`
	Mt          *types.JSON            `json:"mt"`
	Message     string                 `json:"message"`
	BlockHeight uint64                 `json:"block_height"`
	Timestamp   string                 `json:"timestamp"`
	Tag         map[string]interface{} `json:"tag"`
}

type Json struct {
	types.JSON
}
