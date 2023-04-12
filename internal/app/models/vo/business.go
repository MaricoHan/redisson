package vo

import pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/buy"

type BuyRequest struct {
	//Base
	Amount      uint64 `json:"amount" validate:"required"`
	Account     string `json:"account" validate:"required"`
	OrderType   uint32 `json:"order_type"`
	OperationId string `json:"operation_id"`
}

type BatchBuyRequest struct {
	//Base
	List        []*pb.BatchBuyList `json:"list"`
	OperationId string             `json:"operation_id"`
}
