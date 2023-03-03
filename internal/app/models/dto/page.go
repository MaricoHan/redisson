package dto

type Page struct {
	Limit      uint64 `json:"limit" `
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	SortBy     string `json:"sort_by"`
	PageKey    string `json:"next_key"`
	CountTotal string `json:"count_total"`
}

type PageRes struct {
	Limit       uint64 `json:"limit" `
	TotalCount  int64  `json:"total_count"`
	PrevPageKey string `json:"last_key"`
	NextPageKey string `json:"next_key"`
}
