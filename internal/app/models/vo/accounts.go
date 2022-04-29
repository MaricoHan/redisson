package vo

type CreateAccountRequest struct {
	//Base  Base
	OperationID string `json:"operation_id" validate:"required"`
	Count       int64  `json:"count"`
}
