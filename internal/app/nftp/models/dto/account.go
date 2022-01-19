package dto

type CreateAccountP struct {
	Count int64  `json:"count"`
	AppID uint64 `json:"app_id"`
}

type AccountsP struct {
	PageP
	Account string `json:"account"`
	AppID   uint64 `json:"app_id"`
}

type AccountsRes struct {
	PageRes
	Accounts []*Account `json:"accounts"`
}

type Account struct {
	Account string `json:"account"`
	Gas     uint64 `json:"gas"`
}
