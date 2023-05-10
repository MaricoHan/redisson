package evm

type CreateContractCall struct {
	AccessMode  int    `json:"access_mode"`
	OperationId string `json:"operation_id"`
	ProjectID   uint64 `json:"project_id"`
	Code        string `json:"code"`
	Module      string `json:"module"`
	From        string `json:"from"`
	To          string `json:"to"`
	Data        string `json:"data"`
	GasLimit    uint64 `json:"gas_limit"`
	Estimation  uint32 `json:"estimation"`
}

type ShowContractCall struct {
	AccessMode int    `json:"access_mode"`
	ProjectID  uint64 `json:"project_id"`
	Code       string `json:"code"`
	Module     string `json:"module"`
	From       string `json:"from"`
	To         string `json:"to"`
	Data       string `json:"data"`
}

type ShowContractCallRes struct {
	Result string `json:"result"`
}
