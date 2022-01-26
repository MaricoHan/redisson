package dto

type EditNftByIndexP struct {
	Index uint64 `json:"index"`
	Name  string `json:"name"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`

	AppID   uint64 `json:"app_id"`
	ClassId string `json:"class_id"`
	Sender  string `json:"owner"`
}

type EditNftByBatchP struct {
	EditNfts []*EditNft `json:"nfts"`
	AppID    uint64     `json:"app_id"`
	ClassId  string     `json:"class_id"`
	Sender   string     `json:"owner"`
}

type EditNft struct {
	Index uint64 `json:"index" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`
}

type DeleteNftByIndexP struct {
	AppID   uint64 `json:"app_id"`
	ClassId string `json:"class_id"`
	Sender  string `json:"owner"`
	Index   uint64 `json:"index"`
}

type DeleteNftByBatchP struct {
	AppID   uint64   `json:"app_id"`
	ClassId string   `json:"class_id"`
	Sender  string   `json:"owner"`
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
	Timestamp   string `json:"timestamp"`

	AppID uint64 `json:"app_id"`
}

type NftsP struct {
	PageP
	Id      string `json:"id"`
	ClassId string `json:"class_id"`
	Owner   string `json:"owner"`
	TxHash  string `json:"tx_hash"`
	Status  string `json:"status"`
	AppID   uint64 `json:"app_id"`
}

type NftsRes struct {
	PageRes
	Nfts []*Nft `json:"nfts"`
}

type Nft struct {
	Id          string `json:"id"`
	Index       uint64 `json:"index"`
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

type NftClassByIds struct {
	ClassId string `json:"class_id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

type CreateNftsRequest struct {
	AppID     uint64 `json:"app_id"`
	ClassId   string `json:"class_id"`
	Name      string `json:"name"`
	Uri       string `json:"uri"`
	UriHash   string `json:"uri_hash"`
	Data      string `json:"data"`
	Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
}

type NftOperationHistoryByIndexP struct {
	PageP
	ClassID   string `json:"class_id"`
	Index     uint64 `json:"index"`
	Signer    string `json:"signer"`
	Txhash    string `json:"tx_hash"`
	Operation string `json:"operation"`
	AppID     uint64 `json:"app_id"`
}

type BNftOperationHistoryByIndexRes struct {
	PageRes
	OperationRecords []*OperationRecord
}

type OperationRecord struct {
	Txhash    string `json:"tx_hash"`
	Operation string `json:"operation"`
	Signer    string `json:"signer"`
	Recipient string `json:"recipient"`
	Timestamp string `json:"timestamp"`
}
