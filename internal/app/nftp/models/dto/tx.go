package dto

type TxResultByTxHashP struct {
	TaskId    string `json:"task_id"`
	ChainId uint64 `json:"chain_id"`
}

type TxResultByTxHashRes struct {
	Type    string `json:"type"`
	TxHash string  `json:"tx_hash"`
	Status  uint64 `json:"status"`
	Message string `json:"message"`
}
