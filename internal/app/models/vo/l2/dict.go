package l2

type TxType struct {
	Module      uint32 `json:"module"`
	Operation   uint32 `json:"operation"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TxTypesRes struct {
	Data []*TxType `json:"data,omitempty"`
}
