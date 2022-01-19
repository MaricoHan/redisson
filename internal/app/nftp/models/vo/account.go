package vo

type CreateAccountRequest struct {
	Base
	Count int64 `json:"count"`
}
