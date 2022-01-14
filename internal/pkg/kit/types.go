package kit

const (
	// OK request success
	OK = "200"
	// Failed request failed
	Failed = "500"
)

// Response define a struct for http request
type Response struct {
	ErrorResp ErrorResp   `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
