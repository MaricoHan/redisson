package vo

type BatchCreateAccountRequest struct {
	OperationID string `json:"operation_id"`
	Count       uint32 `json:"count"`
}

type CreateAccountRequest struct {
	OperationID string `json:"operation_id"`
	Name        string `json:"name"`
}
