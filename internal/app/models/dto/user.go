package dto

// CreateUsers
//  @Description: 创建用户
type CreateUsers struct {
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	Module     string `json:"module"`
	AccessMode int    `json:"access_mode"`
	Code       string `json:"code"`
	Usertype   uint32 `json:"user_type"`
	Individual struct {
		Name            string `json:"name"`
		Region          int    `json:"region"`
		CertificateType int    `json:"certificate_type"`
		CertificateNum  string `json:"certificate_num"`
		PhoneNum        string `json:"phone_num"`
	} `json:"individual"`
	Enterprise struct {
		Name               string `json:"name"`
		RegistrationRegion int    `json:"registration_region"`
		RegistrationNum    string `json:"registration_num"`
		PhoneNum           string `json:"phone_num"`
		BusinessLicense    string `json:"business_license"`
		Email              string `json:"email"`
	} `json:"enterprise"`
}

// CreateUsersRes
//  @Description: 创建用户返回
type CreateUsersRes struct {
	UserId string `json:"user_id"`
	Did    string `json:"did"`
}

// UpdateUsers
//  @Description: 更新用户
type UpdateUsers struct {
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	Module     string `json:"module"`
	AccessMode int    `json:"access_mode"`
	Code       string `json:"code"`
	UserId     string `json:"user_id"`
	PhoneNum   string `json:"phone_num"`
}
