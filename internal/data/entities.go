package data

import (
	"snap-rq/internal/http"
)

type Collection struct {
	Node[[]Request]
}

type Request struct {
	Node[http.Request]
}

type UserSession struct {
	SelectedRequestId string
	SelectedSessionId string
}
