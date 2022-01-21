package dto

type CreateNftClassP struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	Owner       string `json:"owner"`
	AppID       uint64 `json:"app_id"`
}

type NftClassesP struct {
	PageP
	Id     string `json:"id"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	TxHash string `json:"tx_hash"`
	AppID  uint64 `json:"app_id"`
}

type NftClassesRes struct {
	PageRes
	NftClasses []*NftClass `json:"nft_class"`
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
