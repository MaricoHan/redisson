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
	Recipients  []*pb.Recipient        `json:"recipients"`
	Tag         map[string]interface{} `json:"tag"`
	OperationID string                 `json:"operation_id"`
}
