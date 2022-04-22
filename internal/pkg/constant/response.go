package constant

type Response struct {
	ErrorResp *ErrorResp  `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type ErrorResp struct {
	CodeSpace string `json:"code_space"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}
