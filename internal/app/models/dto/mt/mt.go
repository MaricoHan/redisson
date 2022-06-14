package mt

import (
	pb "gitlab.bianjie.ai/avata/chains/api/pb/mt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
)

type IssueRequest struct {
	Code        string          `json:"code"`
	Module      string          `json:"module"`
	ProjectID   uint64          `json:"project_id"`
	ClassID     string          `json:"class_id"`
	Metadata    string          `json:"metadata"`
	Recipients  []*pb.Recipient `json:"recipients"`
	Tag         string          `json:"tag"`
	OperationID string          `json:"operation_id" validate:"required"`
}
type IssueResponse struct {
	OperationID string `json:"operation_id"`
}

type MintRequest struct {
	Code        string          `json:"code"`
	Module      string          `json:"module"`
	ProjectID   uint64          `json:"project_id"`
	ClassID     string          `json:"class_id"`
	MTID        string          `json:"mt_id"`
	Recipients  []*pb.Recipient `json:"recipients"`
	Tag         string          `json:"tag"`
	OperationID string          `json:"operation_id" validate:"required"`
}
type MintResponse struct {
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
	MtId        string        `json:"mt_id"`         // MT名称
	MtClassId   string        `json:"mt_class_id"`   // 类别ID
	MtClassName string        `json:"mt_class_name"` // 类别名称
	Data        string        `json:"data"`          // 自定义链上元数据
	OwnerCount  uint64        `json:"owner_count"`   // MT 拥有者数量(AVATA平台内)
	IssueData   *pb.IssueData `json:"issue_data"`
	MtCount     uint64        `json:"mt_count"`   // mt 当前流通总量
	MintCount   uint64        `json:"mint_count"` // 发行次数
}

type MTListRequest struct {
	dto.Page
	ProjectID uint64 `json:"project_id"`
	MtId      string `json:"mt_id"`       // MT ID
	MtClassId string `json:"mt_class_id"` // 类别ID
	Issuer    string `json:"issuer"`      // 发行者
	TxHash    string `json:"tx_hash"`     // 交易hash
	Module    string `json:"module"`
	Code      string `json:"code"`
}

type MTListResponse struct {
	dto.PageRes
	Mts []*MT `json:"mts"`
}

type MT struct {
	MtId        string `json:"mt_id"`         // MT ID
	MtClassId   string `json:"mt_class_id"`   // MT 类别 ID
	MtClassName string `json:"mt_class_name"` // MT 类别名称
	Issuer      string `json:"issuer"`        // 发行者
	MtCount     uint64 `json:"mt_count"`      // MT 流通总量
	OwnerCount  uint64 `json:"owner_count"`   // MT 拥有者数量
	Timestamp   string `json:"timestamp"`
}
