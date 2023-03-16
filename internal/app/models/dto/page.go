package dto

type Page struct {
	Limit      uint32 `json:"limit" `
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	SortBy     string `json:"sort_by"`
	PageKey    string `json:"next_key"`
	CountTotal string `json:"count_total"`
}

type PageRes struct {
	PrevPageKey string `json:"prev_page_key"`
	NextPageKey string `json:"next_page_key"`
	Limit       uint32 `json:"limit" `
	TotalCount  uint64 `json:"total_count"`
}
