package evm

import (
	"github.com/volatiletech/sqlboiler/types"
)

type TxResultByTxHash struct {
	OperationId string `json:"operation_id"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
}

type TxResultRes struct {
	Module      uint32      `json:"module"`
	Operation   uint32      `json:"operation"`
	TxHash      string      `json:"tx_hash"`
	Status      uint32      `json:"status"`
	Message     string      `json:"message"`
	BlockHeight uint64      `json:"block_height"`
	Timestamp   string      `json:"timestamp"`
	Nft         *types.JSON `json:"nft,omitempty"`
	Ns          *types.JSON `json:"ns,omitempty"`
}

type TxQueueInfo struct {
	OperationId string `json:"operation_id"`
	ProjectID   uint64 `json:"project_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
}

type TxQueueInfoRes struct {
	QueueTotal       uint64 `json:"queue_total"`
	QueueRequestTime string `json:"queue_request_time"`
	QueueCostTime    uint64 `json:"queue_cost_time"`
	TxQueuePosition  uint64 `json:"tx_queue_position"`
	TxRequestTime    string `json:"tx_request_time"`
	TxCostTime       uint64 `json:"tx_cost_time"`
	TxMessage        string `json:"tx_message"`
}

type Json struct {
	types.JSON
}
