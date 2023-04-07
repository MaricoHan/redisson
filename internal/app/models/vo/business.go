package vo

import pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/buy_v2"

type BuyRequest struct {
	//Base
	Amount      uint64 `json:"amount" validate:"required"`
	Account     string `json:"account" validate:"required"`
	OrderType   uint8  `json:"order_type"`
	OperationId string `json:"operation_id"`
}

type BatchBuyRequest struct {
	//Base
	List        []*pb.BatchBuyList `json:"list"`
	OperationId string             `json:"operation_id"`
}
