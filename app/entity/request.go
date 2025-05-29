package entity

import (
	"fmt"
	"snap-rq/app/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PatchRequest struct {
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	MethodType  *string            `json:"method,omitempty"`
	Url         *string            `json:"url,omitempty"`
	Headers     *map[string]string `json:"headers,omitempty"`
	Body        *string            `json:"body,omitempty"`
}

// Core request
type Request struct {
	Id           string            `json:"id"`
	CollectionID string            `json:"collection_id"` // Foreign key
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	MethodType   string            `json:"method"`
	Url          string            `json:"url"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	ModifiedAt   *time.Time        `json:"modified_at,omitempty"`
	RowPosition  int               `json:"row_position"` // User's logical position of the request
}

func (r *Request) ApplyPatch(patch PatchRequest) {
	now := time.Now()

	if patch.Name != nil {
		r.Name = *patch.Name
	}
	if patch.Description != nil {
		r.Description = *patch.Description
	}
	if patch.MethodType != nil {
		r.MethodType = *patch.MethodType
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

func (r Request) AsHttpRequest() http.HttpRequest {
	return http.HttpRequest{
		Method:  r.MethodType,
		Body:    r.Body,
		URL:     r.Url,
		Headers: r.Headers,
	}
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
		r.MethodType,
		r.Url,
		headersStr,
		strings.TrimSpace(r.Body),
	)
}

func NewRequest(
	collectionID, name, description, method, url string,
	headers map[string]string,
	body string,
	rowPosition int,
) Request {
	requestID := fmt.Sprintf("%s-%s", collectionID, uuid.New().String())
	return Request{
		Id:           requestID,
		CollectionID: collectionID,
		Name:         name,
		Description:  description,
		MethodType:   method,
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
	MethodType  string
	RowPosition int
}

func NewRequestBasicFromRequest(r Request) RequestBasic {
	return RequestBasic{
		Id:          r.Id,
		Url:         r.Url,
		Name:        r.Name,
		MethodType:  r.MethodType,
		RowPosition: r.RowPosition,
	}
}
