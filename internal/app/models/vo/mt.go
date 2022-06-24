package vo

type CreateMTClassRequest struct {
	//Base
	Name        string                 `json:"mt_class_name"`
	Data        string                 `json:"data"`
	Owner       string                 `json:"owner"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}

type TransferMTClassRequest struct {
	//Base
	OperationID string                 `json:"operation_id"`
	Recipient   string                 `json:"recipient"`
	Tag         map[string]interface{} `json:"tag"`
}
