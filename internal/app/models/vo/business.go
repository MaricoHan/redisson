package vo

type BuyRequest struct {
	//Base
	OperationID string `json:"operation_id" validate:"required"`
	Amount      int64  `json:"amount" validate:"required`
	Account     string `json:"account" validate:"required"`
	OrderType   string `json:"order_type"`
	OrderId     string `json:"order_id"`
}
