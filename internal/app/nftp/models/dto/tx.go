package dto

type TxResultByTxHashP struct {
	TaskId     string `json:"task_id"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
}

type TxResultByTxHashRes struct {
	Type    string                 `json:"type"`
	TxHash  string                 `json:"tx_hash"`
	Status  uint64                 `json:"status"`
	ClassID string                 `json:"class_id"`
	NftID   string                 `json:"nft_id"`
	Message string                 `json:"message"`
	Tag     map[string]interface{} `json:"tag"`
}
