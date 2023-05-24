package evm

import (
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/evm/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
)

type NftOperationHistoryByNftId struct {
	dto.Page
	ClassID         string `json:"class_id"`
	NftId           uint64 `json:"nft_id"`
	Signer          string `json:"signer"`
	TxHash          string `json:"tx_hash"`
	Operation       uint32 `json:"operation"`
	ProjectID       uint64 `json:"project_id"`
	ChainID         uint64 `json:"chain_id"`
	PlatFormID      uint64 `json:"plat_form_id"`
	Module          string `json:"module"`
	OperationModule string `json:"operation_module"`
	Code            string `json:"code"`
	AccessMode      int    `json:"access_mode"`
}

type NftOperationHistoryByNftIdRes struct {
	dto.PageRes
	OperationRecords []*OperationRecord `json:"operation_records"`
}

type OperationRecord struct {
	TxHash    string `json:"tx_hash"`
	Operation uint32 `json:"operation"`
	Signer    string `json:"signer"`
	Recipient string `json:"recipient"`
	Timestamp string `json:"timestamp"`
}

type CreateNftClass struct {
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	Uri                  string `json:"uri"`
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
	dto.Page
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
	dto.PageRes
	NftClasses []*NftClass `json:"classes"`
}

type NftClass struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Uri       string `json:"uri"`
	Symbol    string `json:"symbol"`
	Owner     string `json:"owner"`
	TxHash    string `json:"tx_hash"`
	Timestamp string `json:"timestamp"`
}

type NftClassRes struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Uri                  string `json:"uri"`
	Symbol               string `json:"symbol"`
	NftCount             uint64 `json:"nft_count"`
	Owner                string `json:"owner"`
	EditableByOwner      uint32 `json:"editable_by_owner"`
	EditableByClassOwner uint32 `json:"editable_by_class_owner"`
	TxHash               string `json:"tx_hash"`
	Timestamp            string `json:"timestamp"`
}

type NftCount struct {
	Count   uint64 `json:"count"`
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
	dto.PageRes
	Nfts []*NFT `json:"nfts"`
}

type NFT struct {
	Id          uint64 `json:"id"`
	ClassId     string `json:"class_id"`
	ClassName   string `json:"class_name"`
	ClassSymbol string `json:"class_symbol"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Owner       string `json:"owner"`
	Status      int32  `json:"status"`
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
	Status      int32  `json:"status"`
	TxHash      string `json:"tx_hash"`
	Timestamp   string `json:"timestamp"`
}

type Nfts struct {
	dto.Page
	Id         uint64    `json:"id"`
	ClassId    string    `json:"class_id"`
	Owner      string    `json:"owner"`
	TxHash     string    `json:"tx_hash"`
	Status     pb.STATUS `json:"status"`
	ProjectID  uint64    `json:"project_id"`
	ChainID    uint64    `json:"chain_id"`
	PlatFormID uint64    `json:"plat_form_id"`
	Module     string    `json:"module"`
	Code       string    `json:"code"`
	AccessMode int       `json:"access_mode"`
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