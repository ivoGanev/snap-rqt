package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RawHttpRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

type UpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Method      *string `json:"method,omitempty"`
	Url         *string `json:"url,omitempty"`
	Headers     *string `json:"headers,omitempty"`
	Body        *string `json:"body,omitempty"`
}

// Core request entity
type Request struct {
	Id           string     `json:"id"`
	CollectionID string     `json:"collection_id"` // Foreign key
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Method       string     `json:"method"`
	Url          string     `json:"url"`
	Headers      string     `json:"headers,omitempty"`
	Body         string     `json:"body,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	ModifiedAt   *time.Time `json:"modified_at,omitempty"`
	RowPosition  int        `json:"row_position"` // User's logical position of the request
}

func HeadersToString(headers map[string]string) string {
	var b strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&b, "%s: %s\n", k, v)
	}
	return b.String()
}

func StringToHeaders(headers string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(headers, "\n")

	for _, line := range lines {
		// TODO: skip invalid characters, e.g. @

		if line == "" {
			continue // skip empty lines
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}

	// TODO: return errors in case of invalid content
	return result
}

func (r *Request) Update(patch UpdateRequest) {
	now := time.Now()

	if patch.Name != nil {
		r.Name = *patch.Name
	}
	if patch.Description != nil {
		r.Description = *patch.Description
	}
	if patch.Method != nil {
		r.Method = *patch.Method
	}
	if patch.Url != nil {
		r.Url = *patch.Url
	}
	if patch.Headers != nil {
		r.Headers = *patch.Headers
	}
	if patch.Body != nil {
		r.Body = *patch.Body
	}

	r.ModifiedAt = &now
}

func (r Request) AsHttpRequest() RawHttpRequest {
	return RawHttpRequest{
		Method:  r.Method,
		Body:    r.Body,
		URL:     r.Url,
		Headers: StringToHeaders(r.Headers),
	}
}

func (r Request) String() string {
	// Convert headers to a single-line formatted string
	headersStr := "None"
	if len(r.Headers) > 0 {
		var headers []string
		for key, value := range r.Headers {
			headers = append(headers, fmt.Sprintf("%s: %s", key, value))
		}
		headersStr = strings.Join(headers, ", ")
	}

	// Format request as a string
	return fmt.Sprintf(
		"Request { Method: %s  URL: %s  Headers: [%s]  Body: %s }",
		r.Method,
		r.Url,
		headersStr,
		strings.TrimSpace(r.Body),
	)
}

func NewRequest(
	collectionID, name, description, method, url string,
	headers string,
	body string,
	rowPosition int,
) Request {
	requestID := uuid.New().String()
	return Request{
		Id:           requestID,
		CollectionID: collectionID,
		Name:         name,
		Description:  description,
		Method:       method,
		Url:          url,
		Headers:      headers,
		Body:         body,
		CreatedAt:    time.Now(),
		RowPosition:  rowPosition,
	}
}

// Basic Request
type RequestBasic struct {
	Id          string
	Url         string
	Name        string
	Method      string
	RowPosition int
}

func NewRequestBasicFromRequest(r Request) RequestBasic {
	return RequestBasic{
		Id:          r.Id,
		Url:         r.Url,
		Name:        r.Name,
		Method:      r.Method,
		RowPosition: r.RowPosition,
	}
}

func (r RequestBasic) String() string {
	return fmt.Sprintf(
		"RequestBasic { ID: %s  Name: %s  Method: %s  URL: %s  Position: %d }",
		r.Id,
		r.Name,
		r.Method,
		r.Url,
		r.RowPosition,
	)
}
