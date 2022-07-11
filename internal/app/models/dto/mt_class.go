package dto

type MTClassShowRequest struct {
	ProjectID uint64 `json:"project_id"`
	ClassID   string `json:"class_id"`
	Status    string `json:"status"`
	Module    string `json:"module"`
	Code      string `json:"code"`
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
	Page
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
}

type MTClassListResponse struct {
	PageRes
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
	Tag         []byte `json:"tag"`
	Module      string `json:"module"`
	Code        string `json:"code"`
	OperationId string `json:"operation_id"`
}
