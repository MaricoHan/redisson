package vo

type CreateAccountRequest struct {
	Count       uint64 `json:"count" validate:"isdefault=1"`
	OperationID string `json:"operation_id" validate:"required"`
}

type AccountRequest struct {
	PageRequest
	Account string `json:"account" validate:"omitempty"`
}
