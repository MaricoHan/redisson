package types

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Post(url, contentType string, body map[string]interface{}) (*http.Response, error) {
	bodys, _ := json.Marshal(body)
	res, err := http.Post(url, contentType, strings.NewReader(string(bodys)))
	return res, err
}
