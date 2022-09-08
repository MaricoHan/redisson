package vo

type AuthData struct {
	ProjectId  uint64 `json:"project_id"`
	ChainId    uint64 `json:"chain_id"`
	PlatformId uint64 `json:"platform_id"`
	Module     string `json:"module"`
	Code       string `json:"code"`
	AccessMode int    `json:"access_mode"`
}
