package dto

type TransferNftClassByIDP struct {
	ClassID    string `json:"class_id"`
	Owner      string `json:"owner"`
	Recipient  string `json:"recipient"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Tag        []byte `json:"tag"`
}

type TransferNftByNftIdP struct {
	ClassID    string `json:"class_id"`
	Owner      string `json:"owner"`
	NftId      string `json:"nft_id"`
	Recipient  string `json:"recipient"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Tag        []byte `json:"tag"`
}

type TransferNftByBatchP struct {
	ClassID    string       `json:"class_id"`
	Owner      string       `json:"owner"`
	Recipients []*Recipient `json:"recipients"`
	ProjectID  uint64       `json:"project_id"`
	ChainID    uint64       `json:"chain_id"`
	PlatFormID uint64       `json:"plat_form_id"`
}

type Recipient struct {
	NftId     string `json:"nft_id"`
	Recipient string `json:"recipient"`
}
