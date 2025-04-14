package app

import (
	"fmt"
	"snap-rq/app/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RequestsService interface {
	GetRequestListItems() ([]RequestListItem, error)
	SaveRequest(request *Request) error
	GetAllRequests() ([]Request, error)
	GetRequest(id string) (Request, error)
	UpdateRequest(request Request) (Request, error)
}

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
}

type RequestListItem struct {
	Id         string
	Url        string
	Name       string
	MethodType string
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
) Request {
	return Request{
		Id:           uuid.New().String(),
		CollectionID: collectionID,
		Name:         name,
		Description:  description,
		MethodType:   method,
		Url:          url,
		Headers:      headers,
		Body:         body,
		CreatedAt:    time.Now(),
	}
}

func NewRequestSimple(r Request) RequestListItem {
	return RequestListItem{
		Id:         r.Id,
		Url:        r.Url,
		Name:       r.Name,
		MethodType: r.MethodType,
	}
}
