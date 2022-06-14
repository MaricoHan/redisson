package mt

import pb "gitlab.bianjie.ai/avata/chains/api/pb/mt"

type IssueRequest struct {
	Code        string          `json:"code"`
	Module      string          `json:"module"`
	ProjectID   uint64          `json:"project_id"`
	ClassID     string          `json:"class_id"`
	Metadata    string          `json:"metadata"`
	Recipients  []*pb.Recipient `json:"recipients"`
	Tag         string          `json:"tag"`
	OperationID string          `json:"operation_id" validate:"required"`
}
type IssueResponse struct {
	OperationID string `json:"operation_id"`
}

type MintRequest struct {
	Code        string          `json:"code"`
	Module      string          `json:"module"`
	ProjectID   uint64          `json:"project_id"`
	ClassID     string          `json:"class_id"`
	MTID        string          `json:"mt_id"`
	Recipients  []*pb.Recipient `json:"recipients"`
	Tag         string          `json:"tag"`
	OperationID string          `json:"operation_id" validate:"required"`
}
type MintResponse struct {
	OperationID string `json:"operation_id"`
}
