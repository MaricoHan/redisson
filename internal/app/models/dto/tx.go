package dto

type TxResultByTxHash struct {
	TaskId     string `json:"task_id"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
}

type TxResultByTxHashRes struct {
	Type        string                 `json:"type"`
	TxHash      string                 `json:"tx_hash"`
	Status      int32                  `json:"status"`
	ClassID     string                 `json:"class_id"`
	NftID       string                 `json:"nft_id"`
	Message     string                 `json:"message"`
	BlockHeight uint64                 `json:"block_height"`
	Timestamp   string                 `json:"timestamp"`
	Tag         map[string]interface{} `json:"tag"`
}