package dto

type EditNftByIndexP struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
	Data string `json:"data"`

	AppID   uint64 `json:"app_id"`
	ClassId string `json:"class_id"`
	Owner   string `json:"owner"`
	Index   uint64 `json:"index"`
}

type EditNftByBatchP struct {
	Index uint64 `json:"index"`
	Name  string `json:"name"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`
	AppID uint64 `json:"app_id"`
}

type DeleteNftByIndexP struct {
}

type DeleteNftByBatchP struct {
}
type NftByIndexP struct {
	Id          string `json:"id"`
	Index       uint64 `json:"index"`
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
	TimeStamp   string `json:"time_stamp"`

	AppID uint64 `json:"app_id"`
}
