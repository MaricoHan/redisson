package nft

import pb "gitlab.bianjie.ai/avata/chains/api/v2/pb/v2/native/nft"

type CreateNftClassRequest struct {
	// Base
	Name            string `json:"name" validate:"required"`
	ClassId         string `json:"class_id"`
	Symbol          string `json:"symbol"`
	Description     string `json:"description"`
	Uri             string `json:"uri"`
	UriHash         string `json:"uri_hash"`
	Data            string `json:"data"`
	Owner           string `json:"owner" validate:"required"`
	EditableByOwner uint32 `json:"editable_by_owner"`
	OperationID     string `json:"operation_id" validate:"required"`
}

type TransferNftClassByIDRequest struct {
	//Base
	OperationID string `json:"operation_id"`
	Recipient   string `json:"recipient" validate:"required"`
}

type TransferNftByNftIdRequest struct {
	// Base
	OperationID string `json:"operation_id"`
	Recipient   string `json:"recipient" validate:"required"`
}

type CreateNftsRequest struct {
	// Base
	OperationID string `json:"operation_id"`
	Name        string `json:"name" validate:"required"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	// 关闭批量发行
	// Amount    int    `json:"amount"`
	Recipient string `json:"recipient"`
}

type BatchCreateNftsRequest struct {
	// Base
	OperationID string                         `json:"operation_id"`
	Name        string                         `json:"name" validate:"required"`
	Uri         string                         `json:"uri"`
	UriHash     string                         `json:"uri_hash"`
	Data        string                         `json:"data"`
	Recipients  []*pb.NFTBatchCreateRecipients `json:"recipients"`
}

type EditNftByIndexRequest struct {
	Name        string `json:"name" validate:"required"`
	Uri         string `json:"uri"`
	UriHash     string `json:"uri_hash"`
	Data        string `json:"data"`
	OperationID string `json:"operation_id"`
}

type DeleteNftByNftIdRequest struct {
	OperationID string `json:"operation_id"`
}

type BatchTransferRequest struct {
	Data        []*pb.NFTBatchTransferData `json:"data"`
	OperationID string                     `json:"operation_id"`
}

type BatchEditRequest struct {
	Nfts        []*pb.NFTBatchEditData `json:"nfts"`
	OperationID string                 `json:"operation_id"`
}

type BatchDeleteRequest struct {
	Nfts        []*pb.NFTIndex `json:"nfts"`
	OperationID string         `json:"operation_id"`
}
