package dto

type TxResultByTxHashP struct {
	Hash  string `json:"hash"`
	AppID uint64 `json:"app_id"`
}

type TxResultByTxHashRes struct {
	Type    string `json:"type"`
	Status  uint64 `json:"status"`
	Message string `json:"message"`
}
