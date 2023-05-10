package evm

type CreateNftClassRequest struct {
	OperationID          string `json:"operation_id"`
	Name                 string `json:"name" validate:"required"`
	Symbol               string `json:"symbol" validate:"required"`
	Uri                  string `json:"uri"`
	EditableByOwner      uint32 `json:"editable_by_owner"`
	EditableByClassOwner uint32 `json:"editable_by_class_owner"`
	Owner                string `json:"owner" validate:"required"`
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
	Uri         string `json:"uri" validate:"required"`
	UriHash     string `json:"uri_hash"`
	Recipient   string `json:"recipient"`
}

type EditNftByIndexRequest struct {
	Uri         string `json:"uri" validate:"required"`
	UriHash     string `json:"uri_hash"`
	OperationID string `json:"operation_id"`
}

type DeleteNftByNftIdRequest struct {
	OperationID string `json:"operation_id"`
}
