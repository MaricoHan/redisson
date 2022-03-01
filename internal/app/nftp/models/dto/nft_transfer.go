package dto

type TransferNftClassByIDP struct {
	ClassID   string `json:"class_id"`
	Owner     string `json:"owner"`
	Recipient string `json:"recipient"`
	ChainId   uint64 `json:"chain_id"`
}

type TransferNftByNftIdP struct {
	ClassID   string `json:"class_id"`
	Owner     string `json:"owner"`
	NftId     string `json:"nft_id"`
	Recipient string `json:"recipient"`
	ChainId   uint64 `json:"chain_id"`
}

type TransferNftByBatchP struct {
	ClassID    string       `json:"class_id"`
	Owner      string       `json:"owner"`
	Recipients []*Recipient `json:"recipients"`
	ChainId    uint64       `json:"chain_id"`
}

type Recipient struct {
	NftId     string `json:"nft_id"`
	Recipient string `json:"recipient"`
}
