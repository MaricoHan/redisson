package vo

type CreateNftClassRequest struct {
	//Base
	OperationID string                 `json:"operation_id" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Symbol      string                 `json:"symbol"`
	Description string                 `json:"description"`
	Uri         string                 `json:"uri"`
	UriHash     string                 `json:"uri_hash"`
	Data        string                 `json:"data"`
	Owner       string                 `json:"owner" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
}

type TransferNftClassByIDRequest struct {
	//Base
	OperationID string                 `json:"operation_id" validate:"required"`
	Recipient   string                 `json:"recipient" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
}

type TransferNftByNftIdRequest struct {
	//Base
	OperationID string                 `json:"operation_id" validate:"required"`
	Recipient   string                 `json:"recipient" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
}

type CreateNftsRequest struct {
	//Base
	OperationID string `json:"operation_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	//关闭批量发行
	//Amount    int    `json:"amount"`
	Recipient string                 `json:"recipient"`
	Tag       map[string]interface{} `json:"tag"`
}

type EditNftByIndexRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Uri         string                 `json:"uri"`
	Data        string                 `json:"data"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id" validate:"required"`
}

type DeleteNftByNftIdRequest struct {
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id" validate:"required"`
}
