package vo

type CreateMTClassRequest struct {
	//Base
	Name        string                 `json:"name" validate:"required"`
	Data        string                 `json:"data"`
	Owner       string                 `json:"owner" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}
