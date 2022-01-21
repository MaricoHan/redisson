package vo

type CreateNftClassRequest struct {
	Base
	Name        string `json:"name" validate:"required"`
	Symbol      string `json:"symbol" validate:"required"`
	Description string `json:"description" validate:"required"`
	Uri         string `json:"uri" validate:"uri"`
	UriHash     string `json:"uri_hash" validate:"required"`
	Data        string `json:"data" validate:"required"`
	Owner       string `json:"owner" validate:"required"`
}
