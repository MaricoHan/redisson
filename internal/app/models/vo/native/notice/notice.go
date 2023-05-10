package notice

// TransferNFTS 转让NFT通知
type TransferNFTS struct {
	TxHash    string `json:"tx_hash"`
	ProjectID string `json:"project_id"`
}

// TransferClasses 转让Class通知
type TransferClasses struct {
	TxHash    string `json:"tx_hash"`
	ProjectID string `json:"project_id"`
}
