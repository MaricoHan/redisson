package dto

type PageP struct {
	Offset    uint64 `json:"offset"  `
	Limit     uint64 `json:"limit" `
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date" `
	SortBy    string `json:"sort_by"`
}

type PageRes struct {
	Offset     uint64 `json:"offset"  `
	Limit      uint64 `json:"limit" `
	TotalCount uint64 `json:"total_count"`
}
