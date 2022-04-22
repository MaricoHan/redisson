package dto

type Page struct {
	Offset    int64  `json:"offset"`
	Limit     int64  `json:"limit" `
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	SortBy    string `json:"sort_by"`
}

type PageRes struct {
	Offset     int64 `json:"offset"  `
	Limit      int64 `json:"limit" `
	TotalCount int64 `json:"total_count"`
}
