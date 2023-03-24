package dto

type CreateDomain struct {
	OperationId string `json:"operation_id"`
	ProjectID   uint64 `json:"project_id"`
	Code        string `json:"code"`
	Module      string `json:"module"`
	AccessMode  int    `json:"access_mode"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
}

type Domains struct {
	ProjectID  uint64 `json:"project_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
	Name       string `json:"name"`
	Tld        string `json:"tld"`
}

type DomainsRes struct {
	Domains []*Domain `json:"domains"`
}

type Domain struct {
	Name            string `json:"name"`
	Status          uint32 `json:"status"`
	Msg             string `json:"msg"`
	Owner           string `json:"owner"`
	Expire          uint32 `json:"expire"`
	ExpireTimestamp uint64 `json:"expire_timestamp"`
}
