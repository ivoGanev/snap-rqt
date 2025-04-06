package data

import (
	"snap-rq/internal/http"
)

type Collection struct {
	*Node[map[string]Node[http.Request]]
}