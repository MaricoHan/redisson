package dto

import "time"

type PageP struct {
	Offset    int64      `json:"offset"`
	Limit     int64      `json:"limit" `
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	SortBy    string     `json:"sort_by"`
}

type PageRes struct {
	Offset     int64 `json:"offset"  `
	Limit      int64 `json:"limit" `
	TotalCount int64 `json:"total_count"`
}

type TxRes struct {
	TaskId string `json:"task_id"`
}
