package vo

import pb "gitlab.bianjie.ai/avata/chains/api/pb/buy"

type BuyRequest struct {
	//Base
	Amount    uint64 `json:"amount" validate:"required"`
	Account   string `json:"account" validate:"required"`
	OrderType string `json:"order_type"`
	OrderId   string `json:"order_id"`
}

type BatchBuyRequest struct {
	//Base
	List    []*pb.BatchBuyList `json:"list"`
	OrderId string             `json:"order_id"`
}
