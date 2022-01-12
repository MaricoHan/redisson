package kit

const (
	// OK request success
	OK = "200"
	// Failed request failed
	Failed = "500"
)

// Response define a struct for http request
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
