package dto

import (
	pb "gitlab.bianjie.ai/avata/chains/api/pb/mt"
)

type IssueRequest struct {
	Code        string `json:"code"`
	Module      string `json:"module"`
	ProjectID   uint64 `json:"project_id"`
	ClassID     string `json:"class_id"`
	Metadata    string `json:"metadata"`
	Amount      uint64 `json:"amount"`
	Recipient   string `json:"recipient"`
	Tag         string `json:"tag"`
	OperationID string `json:"operation_id" validate:"required"`
}
type IssueResponse struct {
	OperationID string `json:"operation_id"`
}

type MintRequest struct {
	Code        string `json:"code"`
	Module      string `json:"module"`
	ProjectID   uint64 `json:"project_id"`
	ClassID     string `json:"class_id"`
	MTID        string `json:"mt_id"`
	Amount      uint64 `json:"amount,omitempty"`
	Recipient   string `json:"recipient,omitempty"`
	Tag         string `json:"tag"`
	OperationID string `json:"operation_id" validate:"required"`
}
type MintResponse struct {
	OperationID string `json:"operation_id"`
}

type BatchMintRequest struct {
	Code        string          `json:"code"`
	Module      string          `json:"module"`
	ProjectID   uint64          `json:"project_id"`
	ClassID     string          `json:"class_id"`
	MTID        string          `json:"mt_id"`
	Recipients  []*pb.Recipient `json:"recipients"`
	Tag         string          `json:"tag"`
	OperationID string          `json:"operation_id" validate:"required"`
}
type BatchMintResponse struct {
	OperationID string `json:"operation_id"`
}
type EditRequest struct {
	Code        string             `json:"code"`
	Module      string             `json:"module"`
	ProjectID   uint64             `json:"project_id"`
	Owner       string             `json:"owner"`
	Mts         []*pb.EditMetadata `json:"mts"`
	Tag         string             `json:"tag"`
	OperationID string             `json:"operation_id" validate:"required"`
}
type EditResponse struct {
	OperationID string `json:"operation_id"`
}
type BurnRequest struct {
	Code        string `json:"code"`
	Module      string `json:"module"`
	ProjectID   uint64 `json:"project_id"`
	Owner       string `json:"owner"`
	ClassID     string `json:"class_id"`
	MtID        string `json:"mt_id"`
	Amount      uint64 `json:"amount"`
	Tag         string `json:"tag"`
	OperationID string `json:"operation_id" validate:"required"`
}
type BatchBurnRequest struct {
	Code        string       `json:"code"`
	Module      string       `json:"module"`
	ProjectID   uint64       `json:"project_id"`
	Owner       string       `json:"owner"`
	Mts         []*pb.BurnMT `json:"mts"`
	Tag         string       `json:"tag"`
	OperationID string       `json:"operation_id" validate:"required"`
}
type BurnResponse struct {
	OperationID string `json:"operation_id"`
}
type MTBatchTransferRequest struct {
	Code        string         `json:"code"`
	Module      string         `json:"module"`
	ProjectID   uint64         `json:"project_id"`
	Owner       string         `json:"owner"`
	Mts         []*pb.Transfer `json:"mts"`
	Tag         string         `json:"tag"`
	OperationID string         `json:"operation_id" validate:"required"`
}
type BatchBurnResponse struct {
	OperationID string `json:"operation_id"`
}
type MTBatchTransferResponse struct {
	OperationID string `json:"operation_id"`
}

type MTTransferRequest struct {
	Code        string `json:"code"`
	Module      string `json:"module"`
	ProjectID   uint64 `json:"project_id"`
	Owner       string `json:"owner"`
	ClassId     string `json:"class_id"`
	MtId        string `json:"mt_id"`
	Amount      uint64 `json:"amount"`
	Recipient   string `json:"recipient"`
	Tag         string `json:"tag"`
	OperationID string `json:"operation_id" validate:"required"`
}

type MTTransferResponse struct {
	OperationID string `json:"operation_id"`
}

type MTShowRequest struct {
	ProjectID uint64 `json:"project_id"`
	ClassID   string `json:"class_id"`
	MTID      string `json:"mt_id"`
	Module    string `json:"module"`
	Code      string `json:"code"`
}

type MTShowResponse struct {
	MtId       string        `json:"id"`
	ClassId    string        `json:"class_id"`    // 类别ID
	ClassName  string        `json:"class_name"`  // 类别名称
	Data       string        `json:"data"`        // 自定义链上元数据
	OwnerCount uint64        `json:"owner_count"` // MT 拥有者数量(AVATA平台内)
	IssueData  *pb.IssueData `json:"issue_data"`
	MtCount    uint64        `json:"mt_count"`   // mt 当前流通总量
	MintTimes  uint64        `json:"mint_times"` // 发行次数
}

type MTListRequest struct {
	Page
	ProjectID uint64 `json:"project_id"`
	MtId      string `json:"mt_id"`    // MT ID
	ClassId   string `json:"class_id"` // 类别ID
	Issuer    string `json:"issuer"`   // 发行者
	TxHash    string `json:"tx_hash"`  // 交易hash
	Module    string `json:"module"`
	Code      string `json:"code"`
}

type MTListResponse struct {
	PageRes
	Mts []*MT `json:"mts"`
}

type MT struct {
	MtId       string `json:"id"`          // MT ID
	ClassId    string `json:"class_id"`    // MT 类别 ID
	ClassName  string `json:"class_name"`  // MT 类别名称
	Issuer     string `json:"issuer"`      // 发行者
	TxHash     string `json:"tx_hash"`     // MT hash
	OwnerCount uint64 `json:"owner_count"` // MT 拥有者数量
	Timestamp  string `json:"timestamp"`
}

type CreateMTClass struct {
	Name        string `json:"name"`
	Data        string `json:"data"`
	Owner       string `json:"owner"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Tag         []byte `json:"tag"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	OperationId string `json:"operation_id"`
}

type MTOperationHistoryByMTId struct {
	Page
	ClassID    string `json:"class_id"`
	MTId       string `json:"mt_id"`
	Signer     string `json:"signer"`
	Txhash     string `json:"tx_hash"`
	Operation  string `json:"operation"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
}

type MTOperationHistoryByMTIdRes struct {
	PageRes
	OperationRecords []*MTOperationRecord `json:"operation_records"`
}

type MTOperationRecord struct {
	Txhash    string `json:"tx_hash"`
	Operation string `json:"operation"`
	Signer    string `json:"signer"`
	Recipient string `json:"recipient"`
	Amount    uint64 `json:"amount"`
	Timestamp string `json:"timestamp"`
}

type MTBalancesRequest struct {
	Page
	ProjectID uint64 `json:"project_id"`
	MtId      string `json:"mt_id"`    // MT ID
	ClassId   string `json:"class_id"` // 类别ID
	Account   string `json:"account"`
	Module    string `json:"module"`
	Code      string `json:"code"`
}

type MTBalances struct {
	MtId   string `json:"id"` // MT ID
	Amount uint64 `json:"amount"`
}

type MTBalancesResponse struct {
	PageRes
	Mts []*MTBalances `json:"mts"`
}
