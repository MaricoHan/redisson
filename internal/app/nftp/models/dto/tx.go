package dto

type TxResultByTxHashP struct {
	Hash    string `json:"hash"`
	ChainId uint64 `json:"chain_id"`
}

type TxResultByTxHashRes struct {
	Type    string `json:"type"`
	Status  uint64 `json:"status"`
	Message string `json:"message"`
}
