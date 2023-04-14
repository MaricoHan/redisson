package vo

type CreateContractCallRequest struct {
	OperationID string `json:"operation_id"`
	From        string `json:"from" validate:"required"`
	To          string `json:"to" validate:"required"`
	Data        string `json:"data" validate:"required"`
	GasLimit    uint64 `json:"gas_limit" validate:"required"`
	Estimation  uint32 `json:"estimation"`
}
