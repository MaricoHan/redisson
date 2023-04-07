package vo

// CreateUserRequest
//  @Description: 创建用户
type CreateUserRequest struct {
	Usertype           uint32 `json:"user_type"`
	Name               string `json:"name"`
	Region             int    `json:"region"`
	CertificateType    int    `json:"certificate_type"`
	CertificateNum     string `json:"certificate_num"`
	PhoneNum           string `json:"phone_num"`
	RegistrationRegion int    `json:"registration_region"`
	RegistrationNum    string `json:"registration_num"`
	BusinessLicense    string `json:"business_license"`
	Email              string `json:"email"`
}

// UpdateUserRequest
//  @Description: 更新用户
type UpdateUserRequest struct {
	UserId   string `json:"user_id"`
	PhoneNum string `json:"phone_num"`
}

// ShowUserRequest
//  @Description: 查询用户
type ShowUserRequest struct {
	Usertype uint32 `json:"user_type"`
	Code     string `json:"code"`
}
