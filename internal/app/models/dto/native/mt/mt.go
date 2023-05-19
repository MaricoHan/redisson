package mt

import (
	pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
)

type IssueRequest struct {
	Code        string `json:"code"`
	Module      string `json:"module"`
	ProjectID   uint64 `json:"project_id"`
	ClassID     string `json:"class_id"`
	Metadata    string `json:"metadata"`
	Amount      uint64 `json:"amount"`
	Recipient   string `json:"recipient"`
	OperationID string `json:"operation_id" validate:"required"`
	AccessMode  int    `json:"access_mode"`
}
type IssueResponse struct {
}

type MintRequest struct {
	Code        string `json:"code"`
	Module      string `json:"module"`
	ProjectID   uint64 `json:"project_id"`
	ClassID     string `json:"class_id"`
	MTID        string `json:"mt_id"`
	Amount      uint64 `json:"amount,omitempty"`
	Recipient   string `json:"recipient,omitempty"`
	OperationID string `json:"operation_id" validate:"required"`
	AccessMode  int    `json:"access_mode"`
}
type MintResponse struct {
}

type BatchMintRequest struct {
	Code        string          `json:"code"`
	Module      string          `json:"module"`
	ProjectID   uint64          `json:"project_id"`
	ClassID     string          `json:"class_id"`
	MTID        string          `json:"mt_id"`
	Recipients  []*pb.Recipient `json:"recipients"`
	OperationID string          `json:"operation_id" validate:"required"`
	AccessMode  int             `json:"access_mode"`
}
type BatchMintResponse struct {
}
type EditRequest struct {
	Code        string `json:"code"`
	Module      string `json:"module"`
	ProjectID   uint64 `json:"project_id"`
	Owner       string `json:"owner"`
	Data        string `json:"data"`
	ClassId     string `json:"class_id"`
	MTID        string `json:"mt_id"`
	OperationID string `json:"operation_id" validate:"required"`
	AccessMode  int    `json:"access_mode"`
}
type EditResponse struct {
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
	AccessMode  int    `json:"access_mode"`
}
type BatchBurnRequest struct {
	Code        string       `json:"code"`
	Module      string       `json:"module"`
	ProjectID   uint64       `json:"project_id"`
	Owner       string       `json:"owner"`
	Mts         []*pb.BurnMT `json:"mts"`
	OperationID string       `json:"operation_id" validate:"required"`
	AccessMode  int          `json:"access_mode"`
}
type BurnResponse struct {
}
type MTBatchTransferRequest struct {
	Code        string         `json:"code"`
	Module      string         `json:"module"`
	ProjectID   uint64         `json:"project_id"`
	Owner       string         `json:"owner"`
	Mts         []*pb.Transfer `json:"mts"`
	OperationID string         `json:"operation_id" validate:"required"`
	AccessMode  int            `json:"access_mode"`
}
type BatchBurnResponse struct {
}
type MTBatchTransferResponse struct {
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
	OperationID string `json:"operation_id" validate:"required"`
	AccessMode  int    `json:"access_mode"`
}

type MTTransferResponse struct {
}

type MTShowRequest struct {
	ProjectID  uint64 `json:"project_id"`
	ClassID    string `json:"class_id"`
	MTID       string `json:"mt_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
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
	dto.Page
	ProjectID  uint64 `json:"project_id"`
	MtId       string `json:"mt_id"`    // MT ID
	ClassId    string `json:"class_id"` // 类别ID
	Issuer     string `json:"issuer"`   // 发行者
	TxHash     string `json:"tx_hash"`  // 交易hash
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type MTListResponse struct {
	dto.PageRes
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
	Module      string `json:"module"`
	Code        string `json:"code"`
	OperationId string `json:"operation_id"`
	AccessMode  int    `json:"access_mode"`
}

type MTOperationHistoryByMTId struct {
	dto.Page
	ClassID    string `json:"class_id"`
	MTId       string `json:"mt_id"`
	Signer     string `json:"signer"`
	Txhash     string `json:"tx_hash"`
	Operation  uint32 `json:"operation"`
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type MTOperationHistoryByMTIdRes struct {
	dto.PageRes
	OperationRecords []*MTOperationRecord `json:"operation_records"`
}

type MTOperationRecord struct {
	Txhash    string `json:"tx_hash"`
	Operation uint32 `json:"operation"`
	Signer    string `json:"signer"`
	Recipient string `json:"recipient"`
	Amount    uint64 `json:"amount"`
	Timestamp string `json:"timestamp"`
}

type MTBalancesRequest struct {
	dto.Page
	ProjectID  uint64 `json:"project_id"`
	MtId       string `json:"mt_id"`    // MT ID
	ClassId    string `json:"class_id"` // 类别ID
	Account    string `json:"account"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type MTBalances struct {
	MtId   string `json:"id"` // MT ID
	Amount uint64 `json:"amount"`
}

type MTBalancesResponse struct {
	dto.PageRes
	Mts []*MTBalances `json:"mts"`
}

type MTClassShowRequest struct {
	ProjectID  uint64 `json:"project_id"`
	ClassID    string `json:"class_id"`
	Status     string `json:"status"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type MTClassShowResponse struct {
	ClassId   string `json:"id"`        // 类别ID
	ClassName string `json:"name"`      // 类别名称
	Owner     string `json:"owner"`     // 类别所有者
	Data      string `json:"data"`      // 类别拓展数据
	TxHash    string `json:"tx_hash"`   // 交易哈希
	Timestamp string `json:"timestamp"` // 创建时间戳
	MtCount   uint64 `json:"mt_count"`  // MT 类别包含的 MT 总量(AVATA平台内)
}

type MTClassListRequest struct {
	dto.Page
	ProjectID  uint64 `json:"project_id"`
	ClassName  string `json:"name"`    // MT ID
	ClassId    string `json:"id"`      // 类别ID
	Owner      string `json:"owner"`   // 发行者
	TxHash     string `json:"tx_hash"` // 交易hash
	Status     string `json:"status"`
	ChainID    uint64 `json:"chain_id"`
	PlatFormID uint64 `json:"plat_form_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type MTClassListResponse struct {
	dto.PageRes
	MtClasses []*MTClass `json:"classes"`
}

type MTClass struct {
	ClassId   string `json:"id"`
	ClassName string `json:"name"`
	Owner     string `json:"owner"`
	MtCount   uint64 `json:"mt_count"`
	TxHash    string `json:"tx_hash"`
	Timestamp string `json:"timestamp"`
}

type TransferMTClass struct {
	ClassID     string `json:"mt_class_id"`
	Owner       string `json:"owner"`
	Recipient   string `json:"recipient"`
	ProjectID   uint64 `json:"project_id"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	OperationId string `json:"operation_id"`
	AccessMode  int    `json:"access_mode"`
}

type BatchTxRes struct {
}
