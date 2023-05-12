package l2

type CreateNftClassRequest struct {
	OperationID string `json:"operation_id"`
	Name        string `json:"name" validate:"required"`
	ClassId     string `json:"class_id"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	Owner       string `json:"owner" validate:"required"`
}

type TransferNftClassByIDRequest struct {
	// Base
	OperationID string `json:"operation_id"`
	Recipient   string `json:"recipient" validate:"required"`
}

type TransferNftByNftIdRequest struct {
	// Base
	OperationID string `json:"operation_id"`
	Recipient   string `json:"recipient" validate:"required"`
}

type CreateNftsRequest struct {
	OperationID string `json:"operation_id"`
	Name        string `json:"name"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	Recipient   string `json:"recipient"`
}

type EditNftByIndexRequest struct {
	Name        string `json:"name"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	OperationID string `json:"operation_id"`
}

type DeleteNftByNftIdRequest struct {
	OperationID string `json:"operation_id"`
}
