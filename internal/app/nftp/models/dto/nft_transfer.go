package dto

type TransferNftClassByIDP struct {
	ClassID   uint64 `json:"class_id"`
	Owner     string `json:"owner"`
	Recipient string `json:"recipient"`
	AppID     uint64 `json:"app_id"`
}

type TransferNftByIndexP struct {
	ClassID   uint64 `json:"class_id"`
	Owner     string `json:"owner"`
	Index     uint64 `json:"index"`
	Recipient string `json:"recipient"`
	AppID     uint64 `json:"app_id"`
}

type TransferNftByBatchP struct {
	ClassID    uint64       `json:"class_id"`
	Owner      string       `json:"owner"`
	Recipients []*Recipient `json:"recipients"`
	AppID      uint64       `json:"app_id"`
}

type Recipient struct {
	Index     uint64 `json:"index"`
	Recipient string `json:"recipient"`
}
