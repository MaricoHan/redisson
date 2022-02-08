package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func Request(url, contentType, method string, body map[string]interface{}, timeout time.Duration, ctx context.Context) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	bodys, _ := json.Marshal(body)
	request, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(bodys)))
	if err != nil {
		return nil, err
	}
	request.Header.Add("content-type", contentType)
	res, err := client.Do(request)
	if err != nil {
		return res, err
	}
	return res, err
}
