package vo

type CreateDomainRequest struct {
	OperationID string `json:"operation_id"`
	Name        string `json:"name" validate:"required"`
	Owner       string `json:"owner" validate:"required"`
	Duration    uint32 `json:"duration" validate:"required"`
}

type TransferDomainRequest struct {
	OperationID string `json:"operation_id"`
	Recipient   string `json:"recipient" validate:"required"`
}
