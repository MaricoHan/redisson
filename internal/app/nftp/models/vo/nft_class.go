package vo

type CreateNftClassRequest struct {
	Base
	Name        string                 `json:"name" validate:"required"`
	Symbol      string                 `json:"symbol"`
	Description string                 `json:"description"`
	Uri         string                 `json:"uri"`
	UriHash     string                 `json:"uri_hash"`
	Data        string                 `json:"data"`
	Owner       string                 `json:"owner" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
}
