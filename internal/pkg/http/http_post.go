package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func Post(ctx context.Context, url string, body map[string]interface{}) (*http.Response, error) {
	bodys, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(bodys)))
	if err != nil {
		return nil, err
	}
	request.Header.Add("content-type", "application/json")
	return http.DefaultClient.Do(request)
}
