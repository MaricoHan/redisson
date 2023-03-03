package vo

import (
	pb "gitlab.bianjie.ai/avata/chains/api/pb/mt"
)

type IssueRequest struct {
	Name        string `json:"name"`
	Metadata    string `json:"data"`
	Amount      uint64 `json:"amount"`
	Recipient   string `json:"recipient"`
	OperationID string `json:"operation_id"`
}

type MintRequest struct {
	Amount      uint64 `json:"amount,omitempty"`
	Recipient   string `json:"recipient,omitempty"`
	OperationID string `json:"operation_id"`
}
type BatchMintRequest struct {
	Recipients  []*pb.Recipient `json:"recipients"`
	OperationID string          `json:"operation_id"`
}

type EditRequest struct {
	Data        string `json:"data"`
	OperationID string `json:"operation_id"`
}
type BurnRequest struct {
	Amount      uint64 `json:"amount"`
	OperationID string `json:"operation_id"`
}
type BatchBurnRequest struct {
	Mts         []*pb.BurnMT `json:"mts"`
	OperationID string       `json:"operation_id"`
}

type TransferRequest struct {
	Amount      uint64 `json:"amount"`
	Recipient   string `json:"recipient"`
	OperationID string `json:"operation_id"`
}

type BatchTransferRequest struct {
	Mts         []*pb.Transfer `json:"mts"`
	OperationID string         `json:"operation_id"`
}
