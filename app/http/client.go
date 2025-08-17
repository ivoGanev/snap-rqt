package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"snap-rq/app/entity"
	"time"
)

func SendRequest(ctx context.Context, request entity.RawHttpRequest) entity.HttpResult {
	timestamp := time.Now().UnixMilli()
	client := &http.Client{}

	var body io.Reader
	if request.Body != "" {
		body = bytes.NewBuffer([]byte(request.Body))
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, body)
	if err != nil {
		return entity.HttpResult{
			Response: entity.HttpResponse{Timestamp: timestamp},
			Error:    err,
		}
	}

	// Set headers
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return entity.HttpResult{
			Response: entity.HttpResponse{Timestamp: timestamp},
			Error:    err,
		}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.HttpResult{
			Response: entity.HttpResponse{
				Timestamp:  timestamp,
				StatusCode: resp.StatusCode,
			},
			Error: err,
		}
	}

	return entity.HttpResult{
		Response: entity.HttpResponse{
			Timestamp:  timestamp,
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
			Header:     resp.Header,
		},
		Error: nil,
	}
}
