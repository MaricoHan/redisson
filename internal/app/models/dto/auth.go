package dto

type AuthGetUser struct {
	Address   string `json:"address"`
	ChainName string `json:"chain_name"`
}

type AuthVerify struct {
	Exists int `json:"exists"`
}

type AuthUpstreamVerify struct {
	Data struct {
		Exists int `json:"exists"`
	} `json:"data"`
}

// type AuthUpstreamGetUser struct {
// 	Data struct {
// 		Address   string `json:"address"`
// 		ChainName string `json:"chain_name"`
// 	} `json:"data"`
// }

type AuthUpstreamGetUser struct {
	Data []struct {
		Address   string `json:"address"`
		ChainName string `json:"chain_name"`
	} `json:"data"`
}

const (
	AuthVerifyExists = 1
)

var (
	AuthChainName = []interface{}{
		"wenchang-tianhe",
		"wenchang-tianzhou",
	}
)
