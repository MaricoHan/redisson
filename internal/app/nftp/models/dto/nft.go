package dto

type EditNftByIndexP struct {
	Index uint64 `json:"index"`
	Name  string `json:"name"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`

	AppID   uint64 `json:"app_id"`
	ClassId string `json:"class_id"`
	Owner   string `json:"owner"`
}

type EditNftByBatchP struct {
	EditNfts []*EditNft `json:"edit_nfts"`

	AppID   uint64 `json:"app_id"`
	ClassId string `json:"class_id"`
	Owner   string `json:"owner"`
}

type EditNft struct {
	Index uint64 `json:"index"`
	Name  string `json:"name"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`
}

type DeleteNftByIndexP struct {
	AppID   uint64 `json:"app_id"`
	ClassId string `json:"class_id"`
	Owner   string `json:"owner"`
	Index   uint64 `json:"index"`
}

type DeleteNftByBatchP struct {
	AppID   uint64   `json:"app_id"`
	ClassId string   `json:"class_id"`
	Owner   string   `json:"owner"`
	Indices []uint64 `json:"indices"`
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
