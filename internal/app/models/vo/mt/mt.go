package vo

import (
	pb "gitlab.bianjie.ai/avata/chains/api/pb/mt"
)

type IssueRequest struct {
	Name        string                 `json:"name"`
	Metadata    string                 `json:"data"`
	Amount      uint64                 `json:"amount"`
	Recipient   string                 `json:"recipient"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}

type MintRequest struct {
	Amount      uint64                 `json:"amount,omitempty"`
	Recipient   string                 `json:"recipient,omitempty"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}
type BatchMintRequest struct {
	Recipients  []*pb.Recipient        `json:"recipients"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}

type EditRequest struct {
	Data        string                 `json:"data"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}
type BurnRequest struct {
	Amount      uint64                 `json:"amount"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}
type BatchBurnRequest struct {
	Mts         []*pb.BurnMT           `json:"mts"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}

type TransferRequest struct {
	Amount      uint64                 `json:"amount"`
	Recipient   string                 `json:"recipient"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}

type BatchTransferRequest struct {
	Mts         []*pb.Transfer         `json:"mts"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}
