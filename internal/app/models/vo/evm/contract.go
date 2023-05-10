package evm

type CreateContractCallRequest struct {
	OperationID string `json:"operation_id"`
	From        string `json:"from" validate:"required"`
	To          string `json:"to" validate:"required"`
	Data        string `json:"data" validate:"required"`
	GasLimit    uint64 `json:"gas_limit"`
	Estimation  uint32 `json:"estimation"`
}
