package vo

type Base struct {
	OperationID string `json:"operation_id" validate:"required"`
}
