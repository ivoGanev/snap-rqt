package internal

import (
	"snap-rq/internal/data"
	"snap-rq/internal/http"
)

type Collection struct {
	data.Node[[]Request]
}

type Request struct {
	data.Node[http.Request]
}