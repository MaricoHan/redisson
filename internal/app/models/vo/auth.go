package vo

import "encoding/json"

type AuthData struct {
	ProjectId          uint64 `json:"project_id"`
	ChainId            uint64 `json:"chain_id"`
	PlatformId         uint64 `json:"platform_id"`
	Module             string `json:"module"`
	Code               string `json:"code"`
	AccessMode         int    `json:"access_mode"`
	UserId             uint64 `json:"user_id"`
	ExistWalletService bool   `json:"exist_wallet_service"`
}

// AuthVerify auth模块下验证
type AuthVerify struct {
	Hash      string `json:"hash"`
	Type      string `json:"type"`
	ProjectID string `json:"project_id"`
}

// AuthGetUser auth模块下获取账户信息
type AuthGetUser struct {
	Hash      string `json:"hash"`
	Type      string `json:"type"`
	ProjectID string `json:"project_id"`
}

// Map 转换成map
func (a *AuthVerify) Map() (map[string]interface{}, error) {
	var res map[string]interface{}
	msgByte, err := json.Marshal(a)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(msgByte, &res); err != nil {
		return res, err
	}
	return res, nil
}

// Map 转换成map
func (a *AuthGetUser) Map() (map[string]interface{}, error) {
	var res map[string]interface{}
	msgByte, err := json.Marshal(a)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(msgByte, &res); err != nil {
		return res, err
	}
	return res, nil
}
