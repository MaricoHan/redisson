package vo

import "time"

type PageRequest struct {
	Offset    uint64     `json:"offset"  validate:"isdefault=1"`
	Limit     uint64     `json:"limit"  validate:"isdefault=10"`
	StartDate *time.Time `json:"start_date" validate:"datetime"`
	EndDate   *time.Time `json:"end_date" validate:"datetime"`
	SortBy    string     `json:"sort_by"`
}
