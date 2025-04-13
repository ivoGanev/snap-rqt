package http

import (
	"fmt"
	"strings"
)

type RequestMethod string

const (
	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	PUT    RequestMethod = "PUT"
	DELETE RequestMethod = "DELETE"
	PATCH  RequestMethod = "PATCH"
)

var RequestMethods = []RequestMethod{GET, POST, PUT, DELETE, PATCH}

type Request struct {
	Url     string
	Headers map[string]string
	Method  RequestMethod
	Body    string
}



func (r Request) String() string {
	// Convert headers to a formatted string
	headersStr := "None"
	if len(r.Headers) > 0 {
		var headers []string
		for key, value := range r.Headers {
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}
		headersStr = strings.Join(headers, "\n  ")
	}

	// Format request as a string
	return fmt.Sprintf(
		"Request{\n  Method: %s\n  URL: %s\n  Headers:\n  %s\n  Body:\n  %s\n}",
		r.Method,
		r.Url,
		headersStr,
		strings.TrimSpace(r.Body),
	)
}