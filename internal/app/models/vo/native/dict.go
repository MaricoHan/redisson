package native

type TxType struct {
	Module      uint32 `json:"module,omitempty"`
	Operation   uint32 `json:"operation,omitempty"`
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type TxTypesRes struct {
	Data []*TxType `json:"data,omitempty"`
}
