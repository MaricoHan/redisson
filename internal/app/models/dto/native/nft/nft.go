package nft

import (
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/nft"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
)

type TxRes struct {
}

type NftOperationHistoryByNftId struct {
	dto.Page
	ClassID         string `json:"class_id"`
	NftId           string `json:"nft_id"`
	Signer          string `json:"signer"`
	Txhash          string `json:"tx_hash"`
	Operation       string `json:"operation"`
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
	Txhash    string `json:"tx_hash"`
	Operation string `json:"operation"`
	Signer    string `json:"signer"`
	Recipient string `json:"recipient"`
	Timestamp string `json:"timestamp"`
}

type CreateNftClass struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	Owner       string `json:"owner"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
	ClassId     string `json:"class_id"`
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

type BatchTxRes struct {
}

type NftClassesRes struct {
	dto.PageRes
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
	Id          string `json:"id"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	TxHash      string `json:"tx_hash"`
	Symbol      string `json:"symbol"`
	NftCount    uint64 `json:"nft_count"`
	Uri         string `json:"uri"`
	Timestamp   string `json:"timestamp"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	Description string `json:"description"`
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
	NftId       string `json:"nft_id"`
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

type NftReq struct {
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

type Nfts struct {
	dto.Page
	Id         string `json:"id"`
	ClassId    string `json:"class_id"`
	Owner      string `json:"owner"`
	TxHash     string `json:"tx_hash"`
	Status     string `json:"status"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	AccessMode int    `json:"access_mode"`
}

type CreateNfts struct {
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	ClassId     string `json:"class_id"`
	Name        string `json:"name"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	Amount      int    `json:"amount"`
	Recipient   string `json:"recipient"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
}

type BatchCreateNfts struct {
	ProjectID   uint64                         `json:"project_id"`
	ChainID     uint64                         `json:"chain_id"`
	PlatFormID  uint64                         `json:"plat_form_id"`
	ClassId     string                         `json:"class_id"`
	Name        string                         `json:"name"`
	Uri         string                         `json:"uri"`
	UriHash     string                         `json:"uri_hash"`
	Data        string                         `json:"data"`
	Amount      int                            `json:"amount"`
	Recipients  []*pb.NFTBatchCreateRecipients `json:"recipients"`
	Module      string                         `json:"module"`
	Code        string                         `json:"code"`
	OperationId string                         `json:"operation_id"`
	AccessMode  int                            `json:"access_mode"`
}

type NftByNftId struct {
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	NftId      string `json:"nft_id"`
	ClassId    string `json:"class_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type EditNftByNftId struct {
	NftId       string `json:"nft_id"`
	Name        string `json:"name"`
	Uri         string `json:"uri"`
	Data        string `json:"data"`
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
	NftId       string `json:"nft_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
}

type BatchTransferRequest struct {
	Module     string `json:"module"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Sender     string `json:"owner"`
	Code       string `json:"code"`

	Data        []*pb.NFTBatchTransferData `json:"data" validate:"required"`
	OperationID string                     `json:"operation_id" validate:"required"`
	AccessMode  int                        `json:"access_mode"`
}

type BatchEditRequest struct {
	Module     string `json:"module"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Sender     string `json:"owner"`
	Code       string `json:"code"`

	Nfts        []*pb.NFTBatchEditData `json:"nfts"`
	OperationID string                 `json:"operation_id" validate:"required"`
	AccessMode  int                    `json:"access_mode"`
}

type BatchDeleteRequest struct {
	Module     string `json:"module"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Sender     string `json:"owner"`
	Code       string `json:"code"`

	Nfts        []*pb.NFTIndex `json:"nfts"`
	OperationID string         `json:"operation_id" validate:"required"`
	AccessMode  int            `json:"access_mode"`
}
