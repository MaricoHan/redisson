package mt

import (
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
)

type MTClassShowRequest struct {
	ProjectID uint64 `json:"project_id"`
	ClassID   string `json:"class_id"`
	Status    string `json:"status"`
	Module    string `json:"module"`
	Code      string `json:"code"`
}

type MTClassShowResponse struct {
	Id          uint64 `json:"id"`            // 主键ID
	MtClassId   string `json:"mt_class_id"`   // 类别ID
	MtClassName string `json:"mt_class_name"` // 类别名称
	Owner       string `json:"owner"`         // 类别所有者
	Data        string `json:"data"`          // 类别拓展数据
	Status      string `json:"status"`        // 状态
	LockedBy    uint64 `json:"locked_by"`     // 被锁定的交易id
	TxHash      string `json:"tx_hash"`       // 交易哈希
	Timestamp   string `json:"timestamp"`     // 创建时间戳
	MtCount     uint64 `json:"mt_count"`      // MT 类别包含的 MT 总量(AVATA平台内)
	Extra1      string `json:"extra1"`        // 扩展字段1
	Extra2      string `json:"extra2"`        // 扩展字段2
	CreatedAt   string `json:"created_at"`    // 数据存入日期
	UpdatedAt   string `json:"updated_at"`
}

type MTClassListRequest struct {
	dto.Page
	ProjectID   uint64 `json:"project_id"`
	MtClassName string `json:"mt_class_name"` // MT ID
	MtClassId   string `json:"mt_class_id"`   // 类别ID
	Owner       string `json:"owner"`         // 发行者
	TxHash      string `json:"tx_hash"`       // 交易hash
	Status      string `json:"status"`
	ChainID     uint64 `json:"chain_id"`
	PlatFormID  uint64 `json:"plat_form_id"`
	Module      string `json:"module"`
	Code        string `json:"code"`
}

type MTClassListResponse struct {
	dto.PageRes
	MtClasses []*MTClass `json:"mt_classes"`
}

type MTClass struct {
	dto.Page
	MtClassId   string `json:"mt_class_id"`
	MtClassName string `json:"mt_class_name"`
	Owner       string `json:"owner"`
	MtCount     uint64 `json:"mt_count"`
	TxHash      string `json:"tx_hash"`
	Timestamp   string `json:"timestamp"`
}
