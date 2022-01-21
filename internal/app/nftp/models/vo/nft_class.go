package vo

type CreateNftClassRequest struct {
	Base
	Name        string `json:"name" validate:"required"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Uri         string `json:"uri" validate:"uri"`
	UriHash     string `json:"uri_hash" validate:"hexadecimal"`
	Data        string `json:"data"`
	Owner       string `json:"owner" validate:"required"`
}
