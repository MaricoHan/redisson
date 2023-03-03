package vo

type CreateMTClassRequest struct {
	//Base
	Name        string `json:"name"`
	Data        string `json:"data"`
	Owner       string `json:"owner"`
	OperationID string `json:"operation_id"`
}

type TransferMTClassRequest struct {
	//Base
	OperationID string `json:"operation_id"`
	Recipient   string `json:"recipient"`
}
