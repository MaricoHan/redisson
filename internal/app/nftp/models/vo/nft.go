package vo

type EditNftByIndexRequest struct {
	Base
	Name string `json:"name"`
	Uri  string `json:"uri"`
	Data string `json:"data"`
}

type EditNftByBatchRequest struct {
	Base
	Index uint64 `json:"index"`
	Name  string `json:"name"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`
}
