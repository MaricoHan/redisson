package vo

type BatchCreateAccountRequest struct {
	OperationID string `json:"operation_id"`
	Count       int32  `json:"count"`
}

type CreateAccountRequest struct {
	OperationID string `json:"operation_id"`
	Name        string `json:"name"`
}
