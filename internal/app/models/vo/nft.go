package vo

import pb "gitlab.bianjie.ai/avata/chains/api/pb/nft"

type CreateNftClassRequest struct {
	//Base
	OperationID string                 `json:"operation_id"`
	Name        string                 `json:"name" validate:"required"`
	Symbol      string                 `json:"symbol"`
	Description string                 `json:"description"`
	Uri         string                 `json:"uri"`
	UriHash     string                 `json:"uri_hash"`
	Data        string                 `json:"data"`
	Owner       string                 `json:"owner" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
	ClassId     string                 `json:"class_id"`
}

type TransferNftClassByIDRequest struct {
	//Base
	OperationID string                 `json:"operation_id"`
	Recipient   string                 `json:"recipient" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
}

type TransferNftByNftIdRequest struct {
	//Base
	OperationID string                 `json:"operation_id"`
	Recipient   string                 `json:"recipient" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
}

type CreateNftsRequest struct {
	//Base
	OperationID string `json:"operation_id"`
	Name        string `json:"name" validate:"required"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	//关闭批量发行
	//Amount    int    `json:"amount"`
	Recipient string                 `json:"recipient"`
	Tag       map[string]interface{} `json:"tag"`
}

type BatchCreateNftsRequest struct {
	//Base
	OperationID string                         `json:"operation_id"`
	Name        string                         `json:"name" validate:"required"`
	Uri         string                         `json:"uri"`
	UriHash     string                         `json:"uri_hash"`
	Data        string                         `json:"data"`
	Recipients  []*pb.NFTBatchCreateRecipients `json:"recipients"`
	Tag         map[string]interface{}         `json:"tag"`
}

type EditNftByIndexRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Uri         string                 `json:"uri"`
	Data        string                 `json:"data"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}

type DeleteNftByNftIdRequest struct {
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}

type BatchTransferRequest struct {
	Data        []*pb.NFTBatchTransferData `json:"data" validate:"required"`
	Tag         map[string]interface{}     `json:"tag"`
	OperationID string                     `json:"operation_id" validate:"required"`
}

type BatchEditRequest struct {
	Nfts        []*pb.NFTBatchEditData `json:"nfts" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id" validate:"required"`
}

type BatchDeleteRequest struct {
	Nfts        []*pb.NFTIndex         `json:"nfts" validate:"required"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id" validate:"required"`
}
