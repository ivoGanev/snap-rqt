package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func SendRequest(ctx context.Context, request HttpRequest) string {
	var client = &http.Client{}

	var req *http.Request
	var err error

	if request.Method == "GET" || request.Method == "DELETE" {
		req, err = http.NewRequestWithContext(ctx, string(request.Method), request.Url, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, string(request.Method), request.Url, bytes.NewBuffer([]byte(request.Body)))
		for key, value := range request.Headers {
			req.Header.Set(key, value)
		}
	}

	if err != nil {
		return string(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		return string(err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(err.Error())
	}

	return string(respBody)
}