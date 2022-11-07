package dto

import pb "gitlab.bianjie.ai/avata/chains/api/pb/buy"

type BuildOrderInfo struct {
	ProjectID  uint64 `json:"project_id"`
	Address    string `json:"address"`
	Amount     uint64 `json:"amount"`
	ChainId    uint64 `json:"chain_id"`
	Module     string `json:"module"`
	OrderType  string `json:"order_type"`
	OrderId    string `json:"order_id"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type GetOrder struct {
	OrderId    string `json:"order_id"`
	Module     string `json:"module"`
	ProjectID  uint64 `json:"project_id"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type GetAllOrder struct {
	Page
	Module string `json:"module"`
	//OrderType string  `json:"order_type"`
	ProjectId  uint64 `json:"project_id"`
	OrderId    string `json:"order_id"`
	Account    string `json:"account"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	SortBy     string `json:"sort_by"`
	SortRule   string `json:"sort_rule"`
	Status     string `json:"sort_rule"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}

type BuyResponse struct {
	OrderId string `json:"order_id"`
}

type OrderOperationRes struct {
	PageRes
	OrderInfos []*OrderInfo `json:"order_infos"`
}

type OrderInfo struct {
	OrderId    string `json:"order_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Account    string `json:"account"`
	Amount     string `json:"amount"`
	Number     string `json:"number"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	OrderType  string `json:"order_type"`
}

type BatchBuyGas struct {
	ProjectID  uint64             `json:"project_id"`
	ChainId    uint64             `json:"chain_id"`
	Module     string             `json:"module"`
	List       []*pb.BatchBuyList `json:"list"`
	OrderId    string             `json:"order_id"`
	Code       string             `json:"code"`
	AccessMode int                `json:"access_mode"`
}