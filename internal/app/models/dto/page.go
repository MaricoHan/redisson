package dto

type Page struct {
	Limit      int64  `json:"limit" `
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	SortBy     string `json:"sort_by"`
	NextKey    string `json:"next_key"`
	CountTotal string `json:"count_total"`
}

type PageRes struct {
	Limit      int64  `json:"limit" `
	TotalCount int64  `json:"total_count"`
	PreKey     string `json:"last_key"`
	NextKey    string `json:"next_key"`
}
