package dto

type NftOperationHistoryByNftId struct {
	Page
	ClassID         string `json:"class_id"`
	NftId           uint64 `json:"nft_id"`
	Signer          string `json:"signer"`
	Txhash          string `json:"tx_hash"`
	Operation       uint64 `json:"operation"`
	ProjectID       uint64 `json:"project_id"`
	ChainID         uint64 `json:"chain_id"`
	PlatFormID      uint64 `json:"plat_form_id"`
	Module          string `json:"module"`
	OperationModule string `json:"operation_module"`
	Code            string `json:"code"`
	AccessMode      int    `json:"access_mode"`
}

type NftOperationHistoryByNftIdRes struct {
	PageRes
	OperationRecords []*OperationRecord `json:"operation_records"`
}

type OperationRecord struct {
	Txhash    string `json:"tx_hash"`
	Operation uint64 `json:"operation"`
	Signer    string `json:"signer"`
	Recipient string `json:"recipient"`
	Timestamp string `json:"timestamp"`
}

type CreateNftClass struct {
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	Uri                  string `json:"uri"`
	UriHash              string `json:"uri_hash"`
	EditableByOwner      uint32 `json:"editable_by_owner"`
	EditableByClassOwner uint32 `json:"editable_by_class_owner"`
	Owner                string `json:"owner"`
	ProjectID            uint64 `json:"project_id"`
	ChainID              uint64 `json:"chain_id"`
	PlatFormID           uint64 `json:"plat_form_id"`
	Module               string `json:"module"`
	Code                 string `json:"code"`
	AccessMode           int    `json:"access_mode"`
	OperationId          string `json:"operation_id"`
}

type NftClasses struct {
	Page
	Id         string `json:"id"`
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	TxHash     string `json:"tx_hash"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type TxRes struct {
}

type NftClassesRes struct {
	PageRes
	NftClasses []*NftClass `json:"classes"`
}

type NftClass struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	TxHash    string `json:"tx_hash"`
	Symbol    string `json:"symbol"`
	NftCount  uint64 `json:"nft_count"`
	Uri       string `json:"uri"`
	Timestamp string `json:"timestamp"`
}
type NftClassRes struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Owner                string `json:"owner"`
	TxHash               string `json:"tx_hash"`
	Symbol               string `json:"symbol"`
	NftCount             uint64 `json:"nft_count"`
	Uri                  string `json:"uri"`
	Timestamp            string `json:"timestamp"`
	UriHash              string `json:"uri_hash"`
	EditableByOwner      uint32 `json:"editable_by_owner"`
	EditableByClassOwner uint32 `json:"editable_by_class_owner"`
}

type NftCount struct {
	Count   int64  `json:"count"`
	ClassId string `json:"class_id"`
}

type TransferNftClassById struct {
	ClassID     string `json:"class_id"`
	Owner       string `json:"owner"`
	Recipient   string `json:"recipient"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
}

type TransferNftByNftId struct {
	ClassID     string `json:"class_id"`
	Sender      string `json:"owner"`
	NftId       uint64 `json:"nft_id"`
	Recipient   string `json:"recipient"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
}

type NftsRes struct {
	PageRes
	Nfts []*NFT `json:"nfts"`
}

type NFT struct {
	Id          uint64 `json:"id"`
	ClassId     string `json:"class_id"`
	ClassName   string `json:"class_name"`
	ClassSymbol string `json:"class_symbol"`
	Uri         string `json:"uri"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
	TxHash      string `json:"tx_hash"`
	Timestamp   string `json:"timestamp"`
}

type NftRes struct {
	Id          uint64 `json:"id"`
	ClassId     string `json:"class_id"`
	ClassName   string `json:"class_name"`
	ClassSymbol string `json:"class_symbol"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
	TxHash      string `json:"tx_hash"`
	Timestamp   string `json:"timestamp"`
}

type Nfts struct {
	Page
	Id         uint64 `json:"id"`
	ClassId    string `json:"class_id"`
	Owner      string `json:"owner"`
	TxHash     string `json:"tx_hash"`
	Status     string `json:"status"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type CreateNfts struct {
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	ClassId     string `json:"class_id"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Recipient   string `json:"recipient"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
}

type NftByNftId struct {
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	NftId      uint64 `json:"nft_id"`
	ClassId    string `json:"class_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type EditNftByNftId struct {
	NftId       uint64 `json:"nft_id"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Module      string `json:"module"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	ClassId     string `json:"class_id"`
	Sender      string `json:"owner"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
}

type DeleteNftByNftId struct {
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	ClassId     string `json:"class_id"`
	Sender      string `json:"owner"`
	NftId       uint64 `json:"nft_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
}
