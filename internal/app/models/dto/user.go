package dto

import "gitlab.bianjie.ai/avata/chains/api/v2/pb/wallet"

// CreateUsers
//  @Description: 创建用户
type CreateUsers struct {
	ProjectID  uint64              `json:"project_id"`
	ChainID    uint64              `json:"chain_id"`
	Module     string              `json:"module"`
	AccessMode int                 `json:"access_mode"`
	Code       string              `json:"code"`
	Usertype   uint32              `json:"user_type"`
	Individual *wallet.INDIVIDUALS `json:"individual"`
	Enterprise *wallet.ENTERPRISES `json:"enterprise"`
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

// ShowUsers
//  @Description: 查询用户
type ShowUsers struct {
	ProjectID  uint64 `json:"project_id"`
	ChainID    uint64 `json:"chain_id"`
	Module     string `json:"module"`
	AccessMode int    `json:"access_mode"`
	Code       string `json:"code"`
	Usertype   uint32 `json:"user_type"`
	UserCode   string `json:"user_code"`
}

// ShowUsersAccountsRes
//  @Description: 查询用户
type ShowUsersAccountsRes struct {
	PageRes
	Accounts []*ShowUsersAccount `json:"accounts"`
}

// ShowUsersAccount
//  @Description: 钱包查询用户返回
type ShowUsersAccount struct {
	Account
	ReadOnly uint32 `json:"read_only"`
}
