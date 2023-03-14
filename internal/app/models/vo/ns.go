package vo

type CreateDomainRequest struct {
	OperationID string `json:"operation_id"`
	Name        string `json:"name" validate:"required"`
	Owner       string `json:"owner" validate:"required"`
}
