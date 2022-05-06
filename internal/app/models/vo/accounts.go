package vo

type BatchCreateAccountRequest struct {
	OperationID string `json:"operation_id" validate:"required"`
	Count       int64  `json:"count"`
}

type CreateAccountRequest struct {
	OperationID string `json:"operation_id" validate:"required"`
	Name        string `json:"name"`
}
