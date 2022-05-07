package vo

type BuyRequest struct {
	//Base
	Amount    int64  `json:"amount" validate:"required`
	Account   string `json:"account" validate:"required"`
	OrderType string `json:"order_type"`
	OrderId   string `json:"order_id"`
}
