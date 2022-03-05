package dto

type EditNftByNftIdP struct {
	NftId string `json:"nft_id"`
	Name  string `json:"name"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`

	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	ClassId    string `json:"class_id"`
	Sender     string `json:"owner"`
	Tag        []byte `json:"tag"`
}

type EditNftByBatchP struct {
	EditNfts   []*EditNft `json:"nfts"`
	ProjectID  uint64     `json:"project_id"`
	ChainID    uint64     `json:"chain_id"`
	PlatFormID uint64     `json:"plat_form_id"`
	ClassId    string     `json:"class_id"`
	Sender     string     `json:"owner"`
}

type EditNft struct {
	NftId string `json:"nft_id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`
}

type DeleteNftByNftIdP struct {
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	ClassId    string `json:"class_id"`
	Sender     string `json:"owner"`
	NftId      string `json:"nft_id"`
	Tag        []byte `json:"tag"`
}

type DeleteNftByBatchP struct {
	ProjectID  uint64   `json:"project_id"`
	ChainID    uint64   `json:"chain_id"`
	PlatFormID uint64   `json:"plat_form_id"`
	ClassId    string   `json:"class_id"`
	Sender     string   `json:"owner"`
	NftIds     []string `json:"nft_ids"`
}

type NftByNftIdP struct {
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	NftId      string `json:"nft_id"`
	ClassId    string `json:"class_id"`
}
type NftR struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ClassId     string `json:"class_id"`
	ClassName   string `json:"class_name"`
	ClassSymbol string `json:"class_symbol"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
	TxHash      string `json:"tx_hash"`
	Timestamp   string `json:"timestamp"`
}

type NftsP struct {
	PageP
	Id         string `json:"id"`
	ClassId    string `json:"class_id"`
	Owner      string `json:"owner"`
	TxHash     string `json:"tx_hash"`
	Status     string `json:"status"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
}

type NftsRes struct {
	PageRes
	Nfts []*Nft `json:"nfts"`
}

type Nft struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ClassId     string `json:"class_id"`
	ClassName   string `json:"class_name"`
	ClassSymbol string `json:"class_symbol"`
	Uri         string `json:"uri"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
	TxHash      string `json:"tx_hash"`
	Timestamp   string `json:"timestamp"`
}

type NftClassByIds struct {
	ClassId string `json:"class_id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

type CreateNftsP struct {
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	ClassId    string `json:"class_id"`
	Name       string `json:"name"`
	Uri        string `json:"uri"`
	UriHash    string `json:"uri_hash"`
	Data       string `json:"data"`
	Amount     int    `json:"amount"`
	Recipient  string `json:"recipient"`
	Tag        []byte `json:"tag"`
}

type NftOperationHistoryByNftIdP struct {
	PageP
	ClassID    string `json:"class_id"`
	NftId      string `json:"nft_id"`
	Signer     string `json:"signer"`
	Txhash     string `json:"tx_hash"`
	Operation  string `json:"operation"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
}

type BNftOperationHistoryByNftIdRes struct {
	PageRes
	OperationRecords []*OperationRecord `json:"operation_records"`
}

type OperationRecord struct {
	Txhash    string `json:"tx_hash"`
	Operation string `json:"operation"`
	Signer    string `json:"signer"`
	Recipient string `json:"recipient"`
	Timestamp string `json:"timestamp"`
}
